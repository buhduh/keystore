#!/usr/bin/env python
"""
python v2
Helper python script to generate the configuration file in the expected format.
"""

import os 
import json

dirPath = os.path.dirname(os.path.realpath(__file__))

domain = raw_input("FQDN?\n")
port = raw_input("Port?\n")
dName = raw_input("Database Name?\n")
dUser = raw_input("Database User?\n")
dHost = raw_input("Database Host?\n")
dPort = raw_input("Database Port(default is 3306)?\n")
dPW = raw_input("Database Password?\n")
pwEncriptionKey = raw_input("Password encryption key?\n")
assetsLoc = raw_input("Assets location\n")

vals = {}
vals["domain"] = domain
vals["database_name"] = dName
vals ["database_user"] = dUser
vals ["database_host"] = dHost
vals ["database_port"] = dPort
vals ["database_password"] = dPW
vals ["encryption_key"] = pwEncriptionKey
vals ["assets_location"] = assetsLoc
vals ["port"] = port

f = open(dirPath + "/" + "config.json", "w")
json.dump(vals, f, indent=2)
f.close()
