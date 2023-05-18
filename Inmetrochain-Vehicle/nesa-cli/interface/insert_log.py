import string
import random

def generate_random_password():
	possible_chars = list(string.ascii_letters + string.digits + "!@#$%^&*")
	password_length = 10 
	password = []
	for i in range(password_length):
		password.append(random.choice(possible_chars))

	random.shuffle(password)

	return ('Random password: '+''.join(password))


