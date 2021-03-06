# [Kissyface](https://www.youtube.com/watch?v=oZBFOH728pg) 😘

## A text analysis tool for exported Telegram chat logs.

I met a rando on reddit one night who was looking for help with some text analysis. They wanted to surprise their (presumably quite nerdy) partner with a detailed analysis of their Telegram message history. Specifically they wanted histograms of the day of the week and the hour of the day as well as graph of messages sent and received over time.

It sounded adorably nerdy and I had some time to kill on a plane so I wrote this.

### Generating Test Data

This redditor supplied a small bit of sample data (quite rightly not having any interest in sharing their private messages with their SO) so I needed to produce a big dataset of fake data to analyze. I started by downloading an [archive of the collected works of Shakespeare](https://www.gutenberg.org/cache/epub/100/pg100.txt)

[generate_text.py](https://github.com/rsalmond/kissyface/blob/master/generate_text.py) reads the Shakespeare corpus, does some hand wavy cleaning and then spits it out in the Telegram format shared by our friendly redditor.

### Use

As this redditor isn't themselves a programmer I wrote the analysis tool in Go to provide a nice cross platform binary they can use rather than fuss around trying to get Python to run.
