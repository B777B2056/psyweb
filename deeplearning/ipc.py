import sys
import warnings
from ImageTools import *
from models import *
warnings.filterwarnings('ignore')

def plot_images(subject_file:str, save_dir, save_name=None):
    '''
     plot all 5 images needed.
    :param subject_file: 受试者的.set文件（绝对路径）
    :param save_dir: 图片的保存文件夹（这里是static文件夹的路径）
    :param save_name: 保存的图片的命名（str: 例如 LiSi_18537441043.jpg,可选，如果为空则采用默认命名；
                        如果不为空则所有图片统一采用该命名）
    :return: None
    '''
    import os
    raw = mne.io.read_raw_eeglab(subject_file)
    TimeSeries_plot(save_path=os.path.join(save_dir,'eeg','TimeSeries.jpg' if save_name is None else save_name), raw=raw, start_time=0.5,duration=3, dpi=300, show=False)
    C = Pie_plot(raw=raw, save_path=os.path.join(save_dir,'poss','Pie.jpg' if save_name is None else save_name), double_decision=True, dpi=300)
    TSNE_plot(C=C, save_path=os.path.join(save_dir,'tsne','TSNE.jpg' if save_name is None else save_name), raw=raw, fname=None, txt_path='%s.txt'%C, freq_bands=(1,44), dpi=300, plot=False)
    plot_bands_histogram(save=os.path.join(save_dir,'power','Band_histogram.jpg' if save_name is None else save_name), raw=raw, fname=None, txt_path='%s.txt'%C,
                         channels= None, C = C, dpi = 300, show = False)
    BandPSD(C=C,save=os.path.join(save_dir,'topographic','PSD_Topomap.jpg' if save_name is None else save_name), raw=raw, fname=None, start=1, duration=None,
            baseline=None, cmap='jet',dpi=300,show=False)

class IpcWrapper:
    def __init__(self):
        pass

    def loop(self):
        while True:
            content = sys.stdin.readline().strip( )
            if content == "stop":
                break
            phone_number = content.split("/")[-1]
            # 拼接eeg数据路径
            set_file_path = content+".set"
            fdt_file_path = content+".fdt"
            # 跑深度学习模型，生成pdf
            plot_images(set_file_path,"../web/views/images/",phone_number+".jpg")
            # 通知pdf已生成
            sys.stdout.flush()
            sys.stdout.buffer.write((phone_number+"\n").encode('ascii'))
            sys.stdout.flush()
        