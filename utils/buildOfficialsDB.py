import pandas as pd
import numpy as np
import requests
import json
import random
import pyarrow

stateList = ['CO','OR','OH']
stateList = ['CO']

rawZipCodes = pd.read_csv("../data/zip_code_database.csv")
selectStateZips = rawZipCodes[rawZipCodes.state.isin(stateList)]
selectStateZips.zip.shape[0]

zipDF = pd.DataFrame()
for zipcode in selectStateZips.zip:
    print("loading data for zipcode: " + str(zipcode) + "...")
    tempReps = requests.get('http://localhost:3000/api/localreps/google/lookup?address='+str(zipcode))
    tempIngestedJSON = json.loads(tempReps.content)
    tempZipDF = pd.DataFrame(tempIngestedJSON['representatives'])

    if zipDF.shape[0]>0:
        print("zipDF is non-empty. Appending...")
        zipDF = zipDF.append(tempZipDF)
        print(zipDF.shape)
    else:
        print("zipDF is empty, create var")
        zipDF = tempZipDF
        print(zipDF.shape)
zipDF
zipResult = zipDF.loc[:, zipDF.columns != 'index'].set_index('office').drop_duplicates()
zipResult['official_guid'] = ''
for i in range(len(zipResult)):
    zipResult['official_guid'][i] = '%030x' % random.randrange(16**30)

zipResult

zipResult.to_csv('../data/officials_v2.csv')



zipDF = pd.read_csv('../data/zipDF.csv')
zipDF.to_parquet('../data/zipDF.parquet.gzip', compression='gzip')




def loadZipDivisionDB(stateList, inFile, outFile):
    rawZipCodes = pd.read_csv("../data/" + inFile)
    selectStateZips = rawZipCodes[rawZipCodes.state.isin(stateList)]
    zipDivisions = pd.DataFrame()
    for zipcode in selectStateZips.zip:
        print("loading data for zipcode: " + str(zipcode) + "...")
        tempReps = requests.get('http://localhost:3000/api/localreps/lookup?address='+str(zipcode))
        tempIngestedJSON = json.loads(tempReps.content)
        tempZipDF = pd.DataFrame(tempIngestedJSON['representatives'])
        tempZipDF['zipcode'] = zipcode
        tempZipDivisions = tempZipDF[["zipcode","division"]].drop_duplicates()

        if zipDivisions.shape[0]>0:
            print("zipDivisions is non-empty. Appending...")
            zipDivisions = zipDivisions.append(tempZipDivisions)
            print(zipDivisions.shape)
        else:
            print("zipDivisions is empty, create var")
            zipDivisions = tempZipDivisions
            print(zipDivisions.shape)

    zipDivisions.to_csv('../data/' + outFile, index=False)

    return 0

# buld zipcode to division mapping
stateList = ['OH']
inFile = "zip_code_database.csv"
outFile = "zip_divisions_db.csv"

loadZipDivisionDB(stateList, inFile, outFile)

# build User DB
