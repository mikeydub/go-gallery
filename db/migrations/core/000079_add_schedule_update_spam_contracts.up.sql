select cron.schedule('@daily', 'update contracts set is_provider_marked_spam = spam.is_spam from alchemy_spam_contracts spam where contracts.chain = spam.chain and contracts.address = spam.address');
