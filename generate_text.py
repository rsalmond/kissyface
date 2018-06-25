import re
from random import randrange, choice
from datetime import datetime, timedelta

DATE_FORMAT = '%d.%m.%Y %H:%M:%S'

def load_data():
    """ read and do a little clean up on the sharkespeare corpus 
        https://www.gutenberg.org/cache/epub/100/pg100.txt
    """
    with open('pg100.txt') as f:
        lines = f.readlines()

    phrases = []
    for line in lines:
        # all the lines of poetry / dialogue begin with at least two spaces
        if line.startswith('  '):
            # there's some technical whatever in there we dont want
            if len([word for word in ('http', '@') if word in line]) > 0:
                continue
            first_word = line.split()[0]
            # remove the character names for lines of dialogue (eg. FLORIZEL.)
            if first_word.isupper() and first_word.endswith('.'):
                line = line.replace(first_word, '')
            # remove quotes
            phrases.append(line.strip().replace('"',''))

    return phrases

def generate_convo(data):
    """ produce some correctly formatted conversation logs """
    names = ('L [@CunningLinguist](you)', 'Priyanka')

    beginning = 1262304000 # jan 1st 2010
    message_time = datetime.fromtimestamp(beginning)

    for line in data:
        # increment the message time by a random number of seconds (up to
        # around 3 hrs)
        message_time = message_time + timedelta(seconds=randrange(0, 10000))
        # pick a random person to be the sender
        sender = choice(names)
        # format appropriately
        yield '{}, {}: {}'.format(message_time.strftime(DATE_FORMAT), sender, line)

if __name__ == '__main__':
    data = load_data()
    for line in generate_convo(data):
        print(line)
