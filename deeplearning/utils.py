import mne,os
from matplotlib.ticker import Formatter

def CalPxx(raw=None,fname=None,freq_band=(1,44),nfft=None,overlap=None):
    assert not (raw is None) or not (fname is None)
    if nfft is None:
        nfft = int(raw.info['sfreq'])
    if overlap is None:
        overlap = int(nfft/2)
    if ((raw is None) and not (fname is None)):
        try:
            raw = mne.io.read_raw_eeglab(fname, preload=False)
        except:
            print("Error: Read from %s failed!"%os.path.basename(fname))
            raise
    try:
        Pxx,freq = mne.time_frequency.psd_welch(raw,fmin=freq_band[0],fmax=freq_band[1],
                                                    n_fft=nfft,n_overlap=overlap)
        return Pxx, freq
    except:
        print("Error: Psd_Welch failed!")
        raise

def Relative_PSD(Pxx):
    # Pxx:N x channels x frequency
    import numpy as np
    if len(Pxx.shape) == 2:
        Pxx = np.expand_dims(Pxx,axis=0)
    N, chs, _ = Pxx.shape
    for n in range(N):
        for i in range(chs):
            Pxx[n,i, :] = Pxx[n,i, :] / Pxx[n,i, :].sum()
    if N == 1:
        Pxx = Pxx[0,:,:]
    return Pxx


def PsdProbFromTxt(fname, dim=16 * 45):
    import numpy as np
    Pxx = []
    print(fname)
    with open(fname,'r') as f:
        lines = f.readlines()
        for each in lines:
            each = each[:-1].split('\t')
            Pxx.append([float(_) for _ in each])
    return np.array(Pxx)

def AvgSpectrum_channel(Pxx, freq_bands:list=None, channels=None):
    # Pxx:channels x frequency
    import numpy as np
    if len(Pxx.shape) == 2:
        Pxx = Pxx[np.newaxis,:,:]  ## N x channels x frequency
    if freq_bands is None:
        freq_bands = [(1, 4), (4, 8), (8, 13), (13, 30), (30, 45)]
    Pxx_bands = np.zeros((Pxx.shape[0], Pxx.shape[1], len(freq_bands)))  ## N x channels x bands
    for n in range(Pxx.shape[0]):
        for band_i in range(len(freq_bands)):
            Pxx_bands[n,:, band_i] = np.mean(Pxx[n,:,(freq_bands[band_i][0]-1):(freq_bands[band_i][1]-1)], axis=-1)

    return Pxx_bands

def CreatePowerProbTxt(save_path,raw_eeg, C, freq=(1,45)):
    from models import predict,SimpleCNN_16
    file_label = '_C'
    f = open(save_path, 'w')
    for file in os.listdir(raw_eeg):
        if (not file.endswith('.set')) or (not file_label in file):
            continue
        raw = mne.io.read_raw_eeglab(os.path.join(raw_eeg, file))
        if raw.times[-1] < 4:
            continue
        Pxx, _ = CalPxx(raw,freq_band=(1,44))
        Pxx = Relative_PSD(Pxx)
        pred_score = predict(SimpleCNN_16(), raw, 'SimpleCNN_16channels.pth', Pxx)
        Pxx = Pxx.reshape(1, -1)
        while (len(Pxx.shape) > 1):
            Pxx = Pxx[0,:]
        for p in Pxx.tolist():
            f.write(str(p) + '\t')
        f.write(str(pred_score[0 if C == 'HC' else 1]))
        f.write('\n')
    f.close()

def Scale_Prob(Prob,scale=(0,1)):
    MaxVal, MinVal = Prob.max(), Prob.min()
    k = (scale[1] - scale[0]) / (MaxVal - MinVal)
    b = MinVal * (scale[1] - scale[0]) / (MinVal - MaxVal)
    Prob = k * Prob + b
    return Prob

def y_update_scale_value(temp, position, ratio):
    result = temp/ratio
    return "%.2f"%result

class MyFormatter(Formatter):

    def __init__(self, func,ratio=10):
        self.func = func
        self.ratio = ratio
    def __call__(self, x, pos=None):
        return self.func(x, pos,self.ratio)

def ConcateImages(C,save_path,img_root='./temp/',img_size=(5700,6550),imgs=None,pos=None):

     ## pos :(50, 50),(2000, 150),(4250, 50), (2000, 1500),(4000, 1500))
     ## img_size : (5700,6550)
    from PIL import Image
    if imgs is None:
        imgs = ['TimeSeries.jpg', 'Pie_double.jpg', 'TSNE.jpg', 'Band_histogram.jpg', 'PSD_Topomap.jpg','Avg_PSD_Topomap-%s.jpg'%C]
    if pos is None:
        pos = ((50,50),(2000,150),(4250,50),(800,1500),(2700,1500),(4200,1500))
    assert isinstance(imgs, (list, tuple)) and isinstance(pos, (list, tuple)) and len(imgs) == len(pos)
    image = Image.new('RGB', img_size, "#FFFFFF")
    for i in range(len(imgs)):
        im = Image.open(os.path.join(img_root,imgs[i]))
        if i == len(imgs):
            im.resize((1500,5000),resample=Image.ANTIALIAS)
        image.paste(im,pos[i])
    image.save(save_path)

def reshape_array(input_array,target_array):
    import numpy as np
    assert len(input_array.shape) == len(target_array.shape) and input_array.shape[-1] == target_array.shape[-1]
    h1,w1,c1 = input_array.shape
    h2,w2,c2 = target_array.shape
    if h1 < h2:
        inout_array = np.concatenate((np.ones((int((h2-h1)/2),w1,c1))*255,input_array,np.ones((h2-int((h2-h1)/2),w1,c1))*255),axis=0)
    elif h2 < h1:
        input_array = input_array[int((h1-h2)/2):h1-int((h1-h2)/2),:,:]
    h1 = input_array.shape[0]
    if w1 < w2:
        input_array = np.concatenate((np.ones((h1,int((w2-w1)/2),c1))*255, input_array, np.ones((h1,w2-int((w2-w1)/2),c1))*255),axis=1)
    elif w2 < w1:
        input_array = input_array[:,int((w1-w2)/2):w1-int((w1-w2)/2),:]
    return input_array
