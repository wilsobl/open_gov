
import pandas as pd
import numpy as np
import requests
import json
import csv

def saveData(URL, headers, outputFile):
    outputFilePath = '../data/' + outputFile
    response = requests.get(url = URL, headers = headers)
    df = json.loads(response.text)
    members_list = df['results'][0]['members']
    outputFileObj = open(outputFilePath, 'w')
    csvWriter = csv.writer(outputFileObj)
    count=0
    for member in members_list:
        if count == 0:
            header = member.keys()
            csvWriter.writerow(header)
            count += 1
        csvWriter.writerow(member.values()) 
    outputFileObj.close()

    # df['results'][0]['members'][0].values()
    # with open(outputFilePath, 'w') as json_file:
    #     json.dump(df, json_file)

#load in configurations
with open('../data/config.json') as config_file:
    config = json.load(config_file)



# save house members list
URL = 'https://api.propublica.org/congress/v1/116/house/members.json'
headers = {'x-api-key': config['proPublica']}
saveData(URL=URL, headers=headers, outputFile='house_members.csv')

# save senate members list
URL = 'https://api.propublica.org/congress/v1/116/senate/members.json'
saveData(URL=URL, headers=headers, outputFile='senate_members.csv')