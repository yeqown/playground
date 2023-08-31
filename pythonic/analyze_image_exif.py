"""
This is a snippet of code that will analyze the EXIF data of an image.
It will print out the EXIF data in a human readable format.

@File: analyze_image_exif.py
@Author: yeqown
"""

import exifread

from geopy.geocoders import Nominatim
from typing import Dict, Any

geoconverter = Nominatim(user_agent="analyze_image_exif.py")

class ImageExif(object):
    """
    define a class to store the exif data of an image
    """

    filename = None
    datetime = None
    latitude = None # 纬度
    longitude = None # 经度

    def __init__(self, filename:str, tags:Dict[str, Any]):
        if tags is None:
            return
        
        self.filename = filename

        # print(tags)

        if "Image DateTime" in tags:
            self.datetime = tags["Image DateTime"]
        if "GPS GPSLatitude" in tags:
            self.latitude = ImageExif.format_gps(tags["GPS GPSLatitude"])
        if "GPS GPSLongitude" in tags:
            self.longitude = ImageExif.format_gps(tags["GPS GPSLongitude"])

    def __str__(self):
        location = "Unknown"
        if self.latitude is not None and self.longitude is not None:
            query = "{},{}".format(self.latitude, self.longitude)
            location = geoconverter.reverse(query=query)

        return " datetime: %s\n latitude: %s\nlongitude: %s\n location: %s" % (self.datetime, self.latitude, self.longitude, location)

    @staticmethod
    def format_gps(data):
        list_tmp=str(data).replace('[', '').replace(']', '').split(',')
        list=[ele.strip() for ele in list_tmp]
        data_sec = int(list[-1].split('/')[0]) /(int(list[-1].split('/')[1])*3600)# 秒的值
        data_minute = int(list[1])/60
        data_degree = int(list[0])
        result=data_degree + data_minute + data_sec
        return result

def analyze(imageFile: str) -> ImageExif:
    with open(imageFile, "rb") as f:
        tags:Dict[str, Any] = exifread.process_file(f)
        if len(tags) == 0:
            return None
        
        return ImageExif(imageFile, tags)
    
    # analyze ended



def main():
    file = "/Users/yeqown/Downloads/IMG_1779.png"
    exif = analyze(file)

    print("image(%s) analyzing finished" % file)

    if exif is None:
        print("no exif data found")
        return

    print(exif)

if __name__ == "__main__":
    main()