from StringIO import StringIO

from requests import post, get
from requests.auth import HTTPBasicAuth
import uservoice

class APIError(Exception): pass

def convert_to_sprintly(subdomain, ticket_number, item):
    result = {
        'type': 'defect',
        'id': item['id'],
        'title': item['subject'],
        'attachments': [],
    }

    result['description'] = "%s\n\n\nUservoice URL: %s" % (
        item['messages'][-1]['body'],
        'https://%s.uservoice.com/admin/tickets/%s' % (
            subdomain,
            ticket_number
        ))

    if 'tags' in item:
        result['tags'] = ','.join(item['tags'])

    for message in item['messages']:
        result['attachments'].extend([x['url'] for x in message['attachments']])

    return result

def get_item_from_uservoice(client, subdomain, number):
    results = client.get_collection(
        '/api/v1/tickets/search?query=number:%s' % number)
    return convert_to_sprintly(subdomain, number, results[0])

def upload_attachments(base_url, product_id, email, api_key, number, attachments):
    files = {}
    for i, attachment in enumerate(attachments):
        response = get(attachment)
        files['file_%s' % i] = StringIO(response.content)

    result = post('%sproducts/%s/items/%s/attachments.json' % (base_url,
                                                               product_id,
                                                               number),
                  files=files,
                  auth=HTTPBasicAuth(email, api_key))
    if result.status_code != 200:
        raise APIError("Got a %s status code instead of 200." %
                       result.status_code)

def create_item_in_sprintly(base_url, product_id, email, api_key, data):
    attachments = data.pop('attachments', [])

    result = post('%sproducts/%s/items.json' % (base_url, product_id),
                  data=data,
                  auth=HTTPBasicAuth(email, api_key))
    if result.status_code != 200:
        raise APIError("Got a %s status code instead of 200." %
                       result.status_code)
    number = result.json['number']
    print "Created #%s" % number
    print "%s/product/%s/#!/item/%s" % (base_url.replace('/api/', ''),
                                      product_id,
                                      number)

    upload_attachments(base_url, product_id, email, api_key, number, attachments)

if __name__ == '__main__':
    import argparse
    import sys
    import os
    import json
    if not os.path.exists('./settings.json'):
        sys.exit("You must create a settings.json file. " \
                 "Use the example as your guide.")

    with open('./settings.json', 'r') as cfg:
        config = json.load(cfg)

    parser = argparse.ArgumentParser(description='Import stuff from UserVoice.')

    parser.add_argument('ticket_number',help='uservoice ticket number')

    args = parser.parse_args()

    _client = uservoice.Client(config['uservoice']['subdomain'],
                              config['uservoice']['api_key'],
                              config['uservoice']['api_secret'])

                                   config['uservoice']['subdomain'],
                                   args.ticket_number)
    _item = get_item_from_uservoice(_client,

    create_item_in_sprintly(config['sprintly']['base_url'],
                            config['sprintly']['product_id'],
                            config['sprintly']['email'],
                            config['sprintly']['api_key'],
                            _item)
