# from flask import Flask, jsonify
import os
from flask import Flask, jsonify, request
from flask_restx import Api, Resource, fields
import json
#imports from blockchain
import sys
from hfc.fabric import Client as client_fabric
from tornado.platform.asyncio import AnyThreadEventLoopPolicy
import asyncio
import base64
import hashlib
import pickle
#from ecdsa import SigningKey, NIST256p
#from ecdsa.util import sigencode_der, sigdecode_der

app = Flask(__name__)

# Descrição do app
app_infos = dict(version='1.0', title='Nesa Blockchain',
                 description='To access the chaincode using this application, select the post method, click "try out", fill in the json file and run.',
                 contact_email='aroliveira@colaborador.inmetro.gov.br'#,prefix='/eee'
                )

# inicia o swagger app
rest_app = Api(app, **app_infos)


insert_model = rest_app.model('Variáveis usadas no primeiro modelo',
	{
    'args': fields.List(cls_or_instance= fields.String, required = True,
					 description="argumentos da função",
					 help="Ex. meter_id, message, b64sig"),
   'domain': fields.List(cls_or_instance= fields.String, required = True,
           description = "Domínios",
                           help="Ex. [inmetro.br, ptb.de]"),
   'channel_name': fields.String(required=True,
                                  description="Nome do canal"),
   'cc_name': fields.String(required = True,
           description = "Nome do chaincode" ),
   'function': fields.String(required = True,
           description = "Função do chaincode a ser executada" )
     })

# @app.route("/consult")
# def consult():
#   #array_do_usuario = np.array([array_do_usuario])
#   #pred = modelo_carregado.predict(array_do_usuario.reshape(1,-1))
#   return (f"sua solicitação foi predita como: ok", 200)

## Vamos organizar os endpoints por aqui!
# link gerado será: http://127.0.0.1:5000/primeiro_endpoint_swagger

    
consult = rest_app.namespace('', description='Here we will execute the chaincode functions .')
@consult.route("/")
class Teste(Resource):
    @rest_app.expect(insert_model)
    #@rest_app.marshal_with(db_model)
    def post(self):
        args = request.json['args']
        domain = request.json['domain']
        channel_name = request.json['channel_name']
        cc_name = request.json['cc_name']
        function = request.json['function']
        return {
                "Function": "Executed",
                "Result" : insert_functions(args, domain, channel_name, cc_name, function)
               }

def insert_functions(args, domain, channel_name, cc_name, function):
    if function == "insertLog":
        return insert_log(args, domain, channel_name, cc_name, function)
    elif function == "insertJson":
        return insert_json(args, domain, channel_name, cc_name, function)
    elif function == "countHistory":
        return count_history(args, domain, channel_name, cc_name, function)
    else:
        return "This function does not exist"

def insert_json(args, domain, channel_name, cc_name, function):

    path = args[0]
    callpeer = []

    file = open(path)
    data = json.load(file)

    for sensor in range(len(data['signals'])):

        id = data['hardware_id'] + data['signals'][sensor]['sensor_id']

        Sensor_data = { 
            'timestamp_signal': data['timestamp_signal'],  
            'sampling_period_in_sec': data['signals'][sensor]['sampling_period_in_sec'],
            'overall_samples:': data['signals'][sensor]['overall_samples:'],
            'sample_rate_hz': data['signals'][sensor]['sample_rate_hz'],
            'total_of_seconds': data['signals'][sensor]['total_of_seconds'],
            'signal_data' : data['signals'][sensor]['signal_data']
        }

        Sensor_data = json.dumps(Sensor_data, indent = 2).encode('UTF-8')


        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)

        #instantiate the hyperledeger fabric client
        c_hlf = client_fabric(net_profile=("inmetro.br.json"))

        #get access to Fabric as Admin user
        admin = c_hlf.get_user(domain[0], 'Admin')

        for i in domain:
            callpeer.append("peer0." + i)

        #the Fabric Python SDK do not read the channel configuration, we need to add it mannually'''
        c_hlf.new_channel(channel_name)
        asyncio.set_event_loop_policy(AnyThreadEventLoopPolicy())
        #invoke the chaincode to register the meter
        response = loop.run_until_complete(c_hlf.chaincode_invoke(
        requestor=admin,
        channel_name=channel_name,
        peers=callpeer,
        args=[id, Sensor_data],
        cc_name=cc_name,
        fcn=function
        ))

        return "Sucess"

def insert_log(args, domain, channel_name, cc_name, function):

    path = args[0]
    callpeer = []

    file = open(path)
    signal = json.load(file)
    signal = json.dumps(signal)
    log = json.loads(signal)	

    for sensor in range(len(log['Logs'])):

        id = log['device_id'] + log['Logs'][sensor]['sensor_id']

        Sensor_log = { 
            'timestamp_signal': log['timestamp_signal'],  
            'log_details': log['Logs'][sensor]['log_details'],
        }

        Sensor_log = json.dumps(Sensor_log, indent = 2).encode('UTF-8')

        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        #instantiate the hyperledeger fabric client
        c_hlf = client_fabric(net_profile=(domain[0] + ".json"))

        #get access to Fabric as Admin user
        admin = c_hlf.get_user(domain[0], 'Admin')

        for i in domain:
            callpeer.append("peer0." + i)

        #the Fabric Python SDK do not read the channel configuration, we need to add it mannually'''
        c_hlf.new_channel(channel_name)
        asyncio.set_event_loop_policy(AnyThreadEventLoopPolicy())
        #invoke the chaincode to register the meter
        response = loop.run_until_complete(c_hlf.chaincode_invoke(
        requestor = admin,
        channel_name = channel_name,
        peers=callpeer,
        args=[id, Sensor_log],
        cc_name=cc_name,
        fcn=function
        ))

        return "Sucess"

def count_history(args, domain, channel_name, cc_name, function):
    
    id = args[0]
    callpeer = []

    #creates a loop object to manage async transactions
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    #instantiate the hyperledeger fabric client
    c_hlf = client_fabric(net_profile=("inmetro.br.json"))

    #get access to Fabric as Admin user
    admin = c_hlf.get_user(domain[0], 'Admin')

    for i in domain:
        callpeer.append("peer0." + i)

#callpeer = "peer0." + domain
#callpeer = "peer0." + domain

#the Fabric Python SDK do not read the channel configuration, we need to add it mannually
    c_hlf.new_channel(channel_name)

    #invoke the chaincode to register the meter
    response = loop.run_until_complete(c_hlf.chaincode_invoke(
    requestor=admin, 
    channel_name=channel_name, 
    peers=callpeer,
    cc_name=cc_name, 
    fcn=function, 
    args=[id], 
    cc_pattern=None))

    #the signature checking returned... (true or false)
    return("CountHistory:\n", response)

if __name__ == "__main__":
  debug = True # com essa opção como True, ao salvar, o "site" recarrega automaticamente.
  port = int(os.environ.get("PORT", 5000))
  app.run(host='0.0.0.0', port=port, debug=debug)
