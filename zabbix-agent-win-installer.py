#! /usr/bin/env python
__author__ = 'jiasir (Taio Jia) <jiasir@icloud.com>'
-*- coding: utf-8 -*-

import os, sys, re, urllib, wmi
from utils.execute import Execute

cmd = Execute()

c = wmi.wmi()

for s in c.Win32_Service():
	if s.Caption == 'Zabbix Agent':
		sys.exit('zabbix already install')

zabbix_url = 'http://172.20.10.75/zabbix_agent.win.rar'
zabbix_local = 'C:\zabbix_agent.win.rar'
urllib.urlretrieve(zabbix_url, zabbix_local)

rar_url = 'http://172.20.10.75/rar.exe'
rar_local = 'c:\\rar.exe'
urllib.urlretrieve(rar_url, rar_local)
cmd.run_getoutput(rar_local, 'x', '-y', zabbix_local, '-ed', 'C:\\')

s_program_files = os.environ['PROGRAMFILES']
if '(86)' in s_program_files:
	cmd.run_getoutput('C:\\zabbix\\bin\\win64\\zabbix_agentd.exe', '-c', 'C:\\zabbix\\conf\\zabbix_agentd.win.conf', '-i')
else:
	cmd.run_getoutput('C:\\zabbix\\bin\\win32\\zabbix_agentd.exe', '-c', 'C:\\zabbix\\conf\\zabbix_agentd.win.conf', '-i')

conm = cmd.output_to_variable('typeperf.exe -qx | findstr "Network Interface" | findstr "Bytes" | findstr /v "Total" | findstr /v "Loopback"')
f = open('C:\\zabbix\\conf\\zabbix_agentd.win.conf', 'a+')
f.write('\n')

e = 0
for i in range(len(conm)):
	c = re.search('Sent', conm[i])
	if c:
		b = 'PerfCounter ' + '=' + ' eth' + str(e) + '_Out,' + '"''"' + str(conm[i]).strip() + '", 30'
		f.write('%s \n'%b)
		e += 1

e = 0
for i in range(len(conm)):
	c = re.search('received', conm[i])
	if c:
		b = 'PerfCounter ' + '=' + ' eth' + str(e) + '_In,' + str(conm[i]).strip() +'", 30'
		f.write('%s \n' %b)
f.close

f = open('C:\\zabbix\zabbix_agentd.conf', 'r+')
ip = f.read()
ip = ip.replace('192.168.1.100', ipnew)
f.seek(0)
w.write(ip)
f.close

cmd.run_getoutput('net', 'start' '"Zabbix Agent"')
os.remove('C:\\rar.exe')
os.remove('C:\\zabbix_agentd.win.rar')
sys.exit('zabbix install success!')
