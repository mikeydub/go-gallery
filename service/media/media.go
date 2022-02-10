package media

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/everFinance/goar"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/service/rpc"
	"github.com/mikeydub/go-gallery/util"
	"github.com/nfnt/resize"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

var errAlreadyHasMedia = errors.New("token already has preview and thumbnail URLs")

var downloadLock = &sync.Mutex{}

type errUnsupportedURL struct {
	url string
}

type errUnsupportedMediaType struct {
	mediaType persist.MediaType
}

// MakePreviewsForMetadata uses a metadata map to generate media content and cache resized versions of the media content.
func MakePreviewsForMetadata(pCtx context.Context, metadata persist.TokenMetadata, contractAddress persist.Address, tokenID persist.TokenID, turi persist.TokenURI, ipfsClient *shell.Shell, arweaveClient *goar.Client, storageClient *storage.Client) (persist.Media, error) {

	name := fmt.Sprintf("%s-%s", contractAddress, tokenID)

	logrus.Infof("Making previews for %s", name)

	res := persist.Media{}

	imgURL, vURL := findInitialURLs(metadata, name, turi)

	asURI := persist.TokenURI(imgURL)
	switch asURI.Type() {
	case persist.URITypeBase64SVG:
		res.MediaType = persist.MediaTypeSVG
		data, err := rpc.GetDataFromURI(pCtx, asURI, ipfsClient, arweaveClient)
		if err != nil {
			return persist.Media{}, fmt.Errorf("error getting data from base64 svg uri %s: %s", asURI, err)
		}
		res.MediaURL = persist.NullString(data)
		res.ThumbnailURL = res.MediaURL
		return res, nil
	case persist.URITypeSVG:
		res.MediaType = persist.MediaTypeSVG
		res.MediaURL = persist.NullString(asURI.String())
		res.ThumbnailURL = res.MediaURL

		return res, nil
	}

	mediaType, err := downloadAndCache(pCtx, imgURL, name, ipfsClient, arweaveClient)
	if err != nil {
		switch err.(type) {
		case rpc.ErrHTTP:
			if err.(rpc.ErrHTTP).Status == http.StatusNotFound {
				mediaType = persist.MediaTypeInvalid
			} else {
				return persist.Media{}, fmt.Errorf("HTTP error downloading img %s: %s", imgURL, err)
			}
		case *net.DNSError:
			mediaType = persist.MediaTypeInvalid
			logrus.WithError(err).Warnf("DNS error downloading img %s: %s", imgURL, err)
		default:
			return persist.Media{}, fmt.Errorf("error downloading img %s of type %s: %s", imgURL, asURI.Type(), err)
		}
	}
	if vURL != "" {
		logrus.Infof("video url for %s: %s", name, vURL)
		mediaType, err = downloadAndCache(pCtx, vURL, name, ipfsClient, arweaveClient)
		if err != nil {
			switch err.(type) {
			case rpc.ErrHTTP:
				if err.(rpc.ErrHTTP).Status == http.StatusNotFound {
					mediaType = persist.MediaTypeInvalid
				} else {
					return persist.Media{}, fmt.Errorf("HTTP error downloading video %s: %s", vURL, err)
				}
			case *net.DNSError:
				mediaType = persist.MediaTypeInvalid
				logrus.WithError(err).Warnf("DNS error downloading video %s: %s", vURL, err)
			default:
				return persist.Media{}, fmt.Errorf("error downloading video %s of type %s: %s", vURL, asURI.Type(), err)
			}
		}
	}

	vURLAsURI := persist.TokenURI(vURL)
	imgURLAsURI := persist.TokenURI(imgURL)

	if vURLAsURI.Type() == persist.URITypeIPFS {
		vURL = strings.ReplaceAll(vURL, "ipfs://", "https://ipfs.io/ipfs/")
	}
	if imgURLAsURI.Type() == persist.URITypeIPFS {
		imgURL = strings.ReplaceAll(imgURL, "ipfs://", "https://ipfs.io/ipfs/")
	}

	logrus.Infof("media type for %s: %s", name, mediaType)

	switch mediaType {
	case persist.MediaTypeImage:
		res = getImageMedia(pCtx, name, storageClient, vURL, imgURL)
	case persist.MediaTypeVideo, persist.MediaTypeAudio, persist.MediaTypeHTML:
		res = getAuxilaryMedia(pCtx, name, storageClient, vURL, imgURL, mediaType)
	default:
		if vURL != "" {
			logrus.Infof("using vURL for %s: %s", name, vURL)
			res.MediaURL = persist.NullString(vURL)
			if imgURL != "" {
				res.ThumbnailURL = persist.NullString(imgURL)
			}
		} else if imgURL != "" {
			logrus.Infof("using imgURL for %s: %s", name, imgURL)
			res.MediaURL = persist.NullString(imgURL)
		}
	}

	logrus.Infof("media for %s of type %s: %+v", name, mediaType, res)

	return res, nil
}

func getAuxilaryMedia(pCtx context.Context, name string, storageClient *storage.Client, vURL string, imgURL string, mediaType persist.MediaType) persist.Media {
	res := persist.Media{
		MediaType: mediaType,
	}
	imageURL, err := getMediaServingURL(pCtx, viper.GetString("GCLOUD_TOKEN_CONTENT_BUCKET"), fmt.Sprintf("image-%s", name), storageClient)
	if err == nil {
		logrus.Infof("found imageURL for %s: %s", name, imageURL)
		res.ThumbnailURL = persist.NullString(imageURL)
	} else {
		imageURL, err = getMediaServingURL(pCtx, viper.GetString("GCLOUD_TOKEN_CONTENT_BUCKET"), fmt.Sprintf("thumbnail-%s", name), storageClient)
		if err == nil {
			logrus.Infof("found thumbnailURL for %s: %s", name, imageURL)
			res.ThumbnailURL = persist.NullString(imageURL)
		} else {
			logrus.WithError(err).Error("could not get image serving URL")
		}
	}
	if vURL != "" {
		logrus.Infof("using vURL %s: %s", name, vURL)
		res.MediaURL = persist.NullString(vURL)
	} else if imageURL != "" {
		logrus.Infof("using imageURL for %s: %s", name, imageURL)
		res.MediaURL = persist.NullString(imageURL)
	} else if imgURL != "" {
		logrus.Infof("using imgURL for %s: %s", name, imgURL)
		res.MediaURL = persist.NullString(imgURL)
	}
	return res
}

func getImageMedia(pCtx context.Context, name string, storageClient *storage.Client, vURL string, imgURL string) persist.Media {
	res := persist.Media{
		MediaType: persist.MediaTypeImage,
	}
	imageURL, err := getMediaServingURL(pCtx, viper.GetString("GCLOUD_TOKEN_CONTENT_BUCKET"), fmt.Sprintf("image-%s", name), storageClient)
	if err == nil {
		logrus.Infof("found imageURL for %s: %s", name, imageURL)
		res.MediaURL = persist.NullString(imageURL)
	} else {
		if vURL != "" {
			logrus.Infof("using vURL for %s: %s", name, vURL)
			res.MediaURL = persist.NullString(vURL)
			if imgURL != "" {
				res.ThumbnailURL = persist.NullString(imgURL)
			}
		} else if imgURL != "" {
			logrus.Infof("using imgURL for %s: %s", name, imgURL)
			res.MediaURL = persist.NullString(imgURL)
		}
	}
	return res
}

func findInitialURLs(metadata persist.TokenMetadata, name string, turi persist.TokenURI) (string, string) {
	imgURL := ""
	vURL := ""

	if it, ok := util.GetValueFromMapUnsafe(metadata, "animation", util.DefaultSearchDepth).(string); ok {
		logrus.Infof("found initial animation url for %s: %s", name, it)
		vURL = it
	} else if it, ok := util.GetValueFromMapUnsafe(metadata, "video", util.DefaultSearchDepth).(string); ok {
		logrus.Infof("found initial video url for %s: %s", name, it)
		vURL = it
	}

	if it, ok := util.GetValueFromMapUnsafe(metadata, "image", util.DefaultSearchDepth).(string); ok {
		logrus.Infof("found initial image url for %s: %s", name, it)
		imgURL = it
	}

	if imgURL == "" {
		logrus.Infof("no image url found for %s - using token URI %s", name, turi)
		imgURL = turi.String()
	}
	return imgURL, vURL
}

func cacheRawMedia(pCtx context.Context, img []byte, bucket, fileName string) error {

	ctx, cancel := context.WithTimeout(pCtx, 5*time.Second)
	defer cancel()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	client.Bucket(bucket).Object(fileName).Delete(ctx)

	sw := client.Bucket(bucket).Object(fileName).NewWriter(ctx)
	_, err = sw.Write(img)
	if err != nil {
		return err
	}

	if err := sw.Close(); err != nil {
		return err
	}
	return nil
}

func getMediaServingURL(pCtx context.Context, bucketID, objectID string, client *storage.Client) (string, error) {
	objectName := fmt.Sprintf("/gs/%s/%s", bucketID, objectID)

	_, err := client.Bucket(bucketID).Object(objectID).Attrs(pCtx)
	if err != nil {
		return "", fmt.Errorf("error getting attrs for %s: %s", objectName, err)
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketID, objectID), nil
	// key, err := blobstore.BlobKeyForFile(pCtx, objectName)
	// if err != nil {
	// 	return "", fmt.Errorf("error getting blob key for %s: %s", objectName, err)
	// }
	// res, err := appimage.ServingURL(pCtx, key, &appimage.ServingURLOptions{Secure: true})
	// if err != nil {
	// 	return "", fmt.Errorf("error getting serving url for %s: %s", objectName, err)
	// }
	// return res.String(), nil
}

func getVideoURL(pCtx context.Context, bucketID, objectID string, client *storage.Client) (string, error) {
	objectName := fmt.Sprintf("/gs/%s/%s", bucketID, objectID)

	_, err := client.Bucket(bucketID).Object(objectID).Attrs(pCtx)
	if err != nil {
		return "", fmt.Errorf("error getting attrs for %s: %s", objectName, err)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketID, objectID), nil
}

func downloadAndCache(pCtx context.Context, url, name string, ipfsClient *shell.Shell, arweaveClient *goar.Client) (mediaType persist.MediaType, err error) {

	asURI := persist.TokenURI(url)

	if asURI.Type() == persist.URITypeBase64SVG {
		return persist.MediaTypeBase64SVG, nil
	} else if asURI.Type() == persist.URITypeSVG {
		return persist.MediaTypeSVG, nil
	}

	// TODO function that checks url for file extension is img and protocol is https

	downloadLock.Lock()
	defer downloadLock.Unlock()
	bs, err := rpc.GetDataFromURI(pCtx, asURI, ipfsClient, arweaveClient)
	if err != nil {
		return mediaType, err
	}

	logrus.Infof("downloaded %b MB from %s", float64(len(bs))/1024/1024, url)

	buf := bytes.NewBuffer(bs)
	mediaType = persist.SniffMediaType(bs)

	logrus.Infof("sniffed media type for %s: %s", url, mediaType)
	switch mediaType {
	case persist.MediaTypeVideo:

		jp, err := thumbnailVideo(pCtx, bs)
		if err != nil {
			logrus.Infof("error generating thumbnail for %s: %s", url, err)
			return mediaType, fmt.Errorf("error generating thumbnail: %s", err)
		}
		jpg, err := jpeg.Decode(bytes.NewBuffer(jp))
		if err != nil {
			return mediaType, fmt.Errorf("error decoding thumbnail as jpeg: %s", err)
		}
		jpg = resize.Thumbnail(1024, 1024, jpg, resize.NearestNeighbor)
		buf = &bytes.Buffer{}
		err = jpeg.Encode(buf, jpg, nil)
		if err != nil {
			return mediaType, fmt.Errorf("error encoding thumbnail as jpeg: %s", err)
		}

		return persist.MediaTypeVideo, cacheRawMedia(pCtx, buf.Bytes(), viper.GetString("GCLOUD_TOKEN_CONTENT_BUCKET"), fmt.Sprintf("thumbnail-%s", name))
	case persist.MediaTypeGIF:
		return persist.MediaTypeGIF, nil
	case persist.MediaTypeHTML:
		return persist.MediaTypeHTML, nil
	case persist.MediaTypeAudio:
		return persist.MediaTypeAudio, nil
	default:
		_, _, err := image.Decode(buf)
		if err != nil {
			_, pngErr := png.Decode(buf)
			_, jpgErr := jpeg.Decode(buf)
			if pngErr != nil && jpgErr != nil {
				return mediaType, fmt.Errorf("could not decode image as jpg or png %s", err.Error())
			}
		}

		switch asURI.Type() {
		case persist.URITypeIPFS, persist.URITypeArweave:
			return persist.MediaTypeImage, cacheRawMedia(pCtx, buf.Bytes(), viper.GetString("GCLOUD_TOKEN_CONTENT_BUCKET"), fmt.Sprintf("image-%s", name))
		default:
			return persist.MediaTypeImage, nil
		}
	}
}

func thumbnailVideo(pCtx context.Context, vid []byte) ([]byte, error) {
	c := exec.Command("ffmpeg", "-i", "pipe:0", "-ss", "00:00:01.000", "-vframes", "1", "-f", "singlejpeg", "pipe:1")
	in := bytes.NewReader(vid)
	c.Stdin = in
	out := &bytes.Buffer{}
	c.Stdout = out
	err := c.Run()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil

}

func pipeIOForCmd(c *exec.Cmd, input []byte) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 5*1024*1024))
	c.Stdout = buf
	c.Stderr = os.Stderr

	stdin, err := c.StdinPipe()
	if err != nil {
		return nil, err
	}
	err = c.Start()
	if err != nil {
		return nil, err
	}

	stdin.Write(input)

	err = stdin.Close()
	if err != nil {
		return nil, err
	}

	err = c.Wait()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func pipeIOForCmdNoOut(c *exec.Cmd, input []byte) error {

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	stdin, err := c.StdinPipe()
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}

	_, err = stdin.Write(input)
	if err != nil {
		return err
	}

	err = stdin.Close()
	if err != nil {
		return err
	}

	err = c.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (e errUnsupportedURL) Error() string {
	return fmt.Sprintf("unsupported url %s", e.url)
}

func (e errUnsupportedMediaType) Error() string {
	return fmt.Sprintf("unsupported media type %s", e.mediaType)
}
