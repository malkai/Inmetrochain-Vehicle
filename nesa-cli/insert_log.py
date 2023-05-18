import json
from hfc.fabric import Client as client_fabric
import asyncio
from tornado.platform.asyncio import AnyThreadEventLoopPolicy

#For convencion, the first domain should be your admin domain.
domain = ["inmetro.br", "nesa.br"]
channel_name = "nmi-channel"
cc_name = "nesa"
cc_version = "1.0"
callpeer = []

f = open('mock_log.json')
signal = json.load(f)
signal = json.dumps(signal)


def insert(signal):

	log = json.loads(signal)

	for sensor in range(len(log['Logs'])):

		id = log['device_id'] + log['Logs'][sensor]['sensor_id']
		
		Sensor_log = { 
			'timestamp_signal': log['timestamp_signal'],  
			'log_details': log['Logs'][sensor]['log_details'],
		}

		Sensor_log = json.dumps(Sensor_log, indent = 2).encode('UTF-8')
		

		loop = asyncio.get_event_loop()

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
		channel_name='nmi-channel',
		peers=['peer0.inmetro.br'],
		args=[id, Sensor_log],
		cc_name=cc_name,
		fcn='insertLog'
		))
		
		return "Sucess"
		
insert(signal)

