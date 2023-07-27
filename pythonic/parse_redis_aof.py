"""
Try to parse redis aof file and get the command list.
Save the command list to a file, savedFilename can be specified which 
is same as the aof file name by default.
"""

import sys
import os

import hiredis
from argparse import ArgumentParser
# from typing import List, Tuple, Dict, Any, Union, Optional

parser = ArgumentParser(
    prog='RedisAOFParser',
    description='parse redis aof file and save the command list to a file')

# define parameters flags and options
parser.add_argument('-f', '--aof', help='aof file name', required=True)
parser.add_argument('-s', '--output', help='saved file name, default is same as aof file name')

def parse_aof(aof_filename:str, saved_filename:str):
    print('parsing aof file: %s' % aof_filename)
    print('saving parsed file: %s' % saved_filename)

    # check if aof file exists
    if not os.path.exists(aof_filename):
        print('aof file not exists')
        sys.exit(1)
    
    # open and read aof file
    reader = hiredis.Reader(encoding='utf-8', errors='strict')
    count:int = 0
    try:
        save_file = os.open(saved_filename, os.O_CREAT|os.O_TRUNC|os.O_RDWR)
    except OSError as e:
        print('open file error: %s' % e)
        sys.exit(1)

    with open(aof_filename, 'rb') as f:
        while True:
            line:str = f.readline()
            if not line:
                # end of file
                break
            reader.feed(line)
            try:
                command = reader.gets()
                if command is not False:
                    line = str(command) + '\n'
                    os.write(save_file, line.encode('utf-8'))
            except hiredis.ProtocolError as e:
                print('protocol error: %s' % e)
                continue
            count += 1

        print("total commands: %d" % count)

    os.close(save_file)

def main():
    args = parser.parse_args()
    aof_filename:str = args.aof
    saved_filename:str = args.output

    if not aof_filename:
        print('aof file name is required')
        sys.exit(1)

    if not saved_filename:
        saved_filename = "parsed_aof.txt"
    
    parse_aof(aof_filename, saved_filename)

if __name__ == '__main__':
    main()