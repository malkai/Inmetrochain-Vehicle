import json
import math
from datetime import datetime, timedelta
import matplotlib.pyplot as plt
import numpy as np
'''

func KalmanFilter(capacicapacidadedade float64, medições []float64) float64 {

	var media = mediavector(medições)

	Gerro := errovector(medições, media) // desvio global

	auxe := []float64{}
	auxe = append(auxe, medições[0])
	Lerro := errovector(auxe, media) // desvio local

	var estima float64
	leiturasPercentuaispos := []float64{}

	for _, leitura := range medições {
		var k = math.Pow(Lerro, 2)/math.Pow(Lerro, 2) + math.Pow(Gerro, 2)
		estima = media + k*(leitura-media)

		Lerro = (1 - k) * Lerro
		leiturasPercentuaispos = append(leiturasPercentuaispos, estima)

	}

	resultadotanque := ((medições[0] - estima) * capacidade) / 100

	return resultadotanque

}

func mediavector(a []float64) float64 {
	media := 0.0
	for _, leitura := range a {
		media = media + leitura
	}

	media = media / float64(len(a))
	return media

}

func errovector(a []float64, m float64) float64 {
	var aux float64
	errovec := []float64{}
	for _, leitura := range a {
		errovec = append(errovec, leitura-m)
	}
	for _, leitura := range errovec {
		aux = +leitura
	}
	aux = aux / float64(len(a))
	return aux
}
'''

def Kalman(medições):
    media = mediavector(medições)

    Gerro =  errovector(medições, media)

    auxe = []
    auxe.append(medições[0])

    Lerro = errovector(auxe, media)

    estima = 0.0
    k = 0.0
    leiturasPercentuaispos = []
  
    for  leitura in  medições:

        k = math.pow(Lerro,2)/(math.pow(Lerro,2) + math.pow(Gerro,2))
        estima = media + k*( float(leitura)-media)
        Lerro = (1-k)*Lerro
        leiturasPercentuaispos.append( round(estima,2))
    resultadotanque = ( round(float(leiturasPercentuaispos[0]),2) - round(estima,2)) 
    
    return resultadotanque,leiturasPercentuaispos


def mediavector(medições):
    media = 0.0
    for  leitura in medições:

        media = media + float(leitura)

    media = media / len(medições)
   
    return media

def errovector(medições,media):
    aux = 0.0
    errovec = []
    for leitura in medições:
        errovec.append(float(leitura)-media)
  
    for leitura in errovec:
        aux = aux+float(leitura)
	
    aux = aux / len(medições)
    
    return aux

def plotlist(comblista, comblistaKalman, vin,timelist):
    y = range(0, len(comblista))
    # the smooth signal
    plt.ioff()
    #print(comblista)
    #print(comblistaKalman)
    y1 =  range(0, len(comblistaKalman))# the raw signal

    fig = plt.figure()
    ax = fig.add_subplot(121)
    ax.plot(y, comblista,label='raw-signal')
    ax.legend(loc='best')

    ax2 = fig.add_subplot(122)
    ax2.plot(y1,comblistaKalman,label='Kalman-signal')
    ax2.legend(loc='best')

    plt.suptitle(vin)
    plt.savefig(vin+'.png')
    
 
# Opening JSON file
f = open('json_data.json')
 
# returns JSON object as 
# a dictionary
data = json.load(f)
f.close()
# Iterating through the json
# list
#b = list(dict.items())
#print(type(data))
#a =  dict(sorted(data.items(), key=lambda item: data[item]))

tuplas = [] 

for i in data:
      #print(data[i])
     
        time = data[i]['timeStamp']
        if(time == ''):
            time = 'NaN' 
        try:  
            fuel =data[i]['obdData']['09 02 5']['response']
            
        except: 
            fuel = 'NaN'
        if(fuel == ''):
            fuel = 'NaN'      
      
        try:
            vin = data[i]['obdData']['01 2F']['response']
        except: 
            vin = 'NaN'
        if(vin == ''):
            vin = 'NaN'  
            
        tuplas.append([time,fuel,vin])
    
lst = sorted(tuplas, key=lambda  x: (x[1], x[0]))


vin = ''
fuelist= []
fuelest = 0.0
fuellistkalman = []
time = 0
timelist = []
for i in lst:
    j = 0
    if lst[j]!=lst[-1]:
        s1 = datetime.strptime(lst[j][0],'%Y-%m-%d %H:%M:%S.%f')
        s2 = datetime.strptime(lst[j+1][0],'%Y-%m-%d %H:%M:%S.%f')
        teste = (s2-s1).seconds 
        time = time + teste
       
        #print(i[2])
        if (vin != i[1]  and vin!='NaN' ) or teste>=600 :
           
            if(len(fuelist)>1):
                fuelest,fuellistkalman = Kalman(fuelist)
            print(vin)
            print(time)
            #print(fuellistkalman)
            plotlist(fuelist,fuellistkalman, vin,timelist)
            vin = i[1]
            fuelist.clear
            timelist.clear
            time = 0
        if(i[2]!= 'NaN' and i[2]!= ' '):
            timelist.append(time)
            fuelist.append(round(float(i[2]),2))   
            
    j=j+1   
    
    

     
'''
Fazer a inserção de um ruído branco nos dados simulados. 
Aplicar o filtro de kalman no ruído gerado. 
Testar com dados previsiveis, depois avaliar os pros e contra dos dadoa reias. 
'''