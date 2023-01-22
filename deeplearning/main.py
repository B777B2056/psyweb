from ipc import *
import os

os.chdir("./deeplearning") 

if __name__ == '__main__':
    ipc = IpcWrapper()
    ipc.loop()
    # content = "/mnt/e/jr/go/src/psyWeb/deeplearning/eeg_data/15901267537"
    # phone_number = content.split("/")[-1]
    # set_file_path = content+".set"
    # fdt_file_path = content+".fdt"
    # plot_images(set_file_path,"../web/views/images/",phone_number+".jpg")
