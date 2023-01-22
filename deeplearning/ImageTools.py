import mne, os, warnings
import matplotlib.pyplot as plt
from matplotlib.backends.backend_agg import FigureCanvasAgg
warnings.filterwarnings('ignore')

def TimeSeries_plot(save_path, fname=None, raw=None, start_time=0.5, duration=4, dpi=300, show=False):
    from PIL import Image,ImageDraw,ImageFont
    import numpy as np
    import matplotlib.font_manager as fm
    assert (fname != None or raw != None)
    if fname:
        try:
            raw = mne.io.read_raw_eeglab(fname, preload=False)
        except:
            print('Read from %s failed!'%os.path.basename(fname))
            return
    try:
        picks = list(range(raw.info['nchan']))
        if raw.info['nchan'] > 16:
            picks = [raw.ch_names.index[_] for _ in ['Fp1', 'Fp2', 'F3', 'F4', 'C3', 'C4', 'P3', 'P4', 'O1',
                                                            'O2', 'F7', 'F8', 'T7', 'T8', 'P7', 'P8']]
        raw.plot(start=start_time, duration=duration, show=show, bgcolor='w',order=picks,n_channels=len(picks))
        plt.tight_layout()
        fig = plt.gcf()
        fig.set_size_inches(20 / 3, 12 / 3)
        plt.savefig(save_path, dpi=dpi)
    except:
        print("Plot image failed!")
        raise

    img = Image.open(save_path)
    img = np.array(img)
    #img = img[:-240,:-100,:]
    img_new = np.ones((img.shape[0] + 100,img.shape[1],img.shape[2]),dtype=np.uint8)*255
    img_new[100:,:,:] = img
    img_pil = Image.fromarray(img_new)
    ttfront = ImageFont.truetype(fm.findfont(fm.FontProperties(family='DejaVu Sans')),90)  # 字体大小
    draw = ImageDraw.Draw(img_pil)
    draw.text(xy=(720,40),text="EEG Time Series",font=ttfront,fill=(0,25,25))
    #ttfront = ImageFont.truetype(fm.findfont(fm.FontProperties(family='DejaVu Sans')),60)
    #ttfront = ImageFont.truetype('arial.ttf', 60)  # 字体大小
    #draw.text(xy=(900, 1160), text="Times (s)", font=ttfront, fill=(0, 25, 25))
    img_pil.save(save_path)

def TSNE_plot(C,save_path,raw=None,fname=None,txt_path= 'HC.txt',freq_bands=(1,44),dpi=300,plot=False):
    from utils import CalPxx,Relative_PSD,PsdProbFromTxt,Scale_Prob
    from sklearn.manifold import TSNE
    import numpy as np
    import os
    assert raw or fname

    if not (raw is None):
        Pxx_cur,freq = CalPxx(raw=raw,freq_band=freq_bands)
    else:
        Pxx_cur, freq = CalPxx(fname=fname, freq_band=freq_bands)
    n_chans = raw.info['nchan']
    assert n_chans in [8,16,32]
    if n_chans == 32:
        selected_channels = [raw.ch_names.index(_) for _ in
                            ['Fp1', 'Fp2', 'F3', 'F4', 'C3', 'C4', 'P3', 'P4', 'O1', 'O2', 'F7', 'F8', 'T7', 'T8', 'P7', 'P8']]
        Pxx_cur = Pxx_cur.squeeze()[selected_channels,:]

    Pxx_cur = Relative_PSD(Pxx_cur)
    Pxx_cur = Pxx_cur.reshape(1,-1)
    Pxx_pre = PsdProbFromTxt(fname=os.path.join('.',txt_path.replace('.','_16channels.') if n_chans >= 16 else txt_path.replace('.','_8channels.')))
    Prob_pre = Pxx_pre[:,-1]
    Pxx_pre = Pxx_pre[:,:-1]
    if n_chans < 16:
        Pxx_pre = Pxx_pre.reshape(Pxx_pre.shape[0],freq_bands[-1],16)
        Pxx_pre = Pxx_pre[:,:,[0,1,2,3,4,5,10,11]].reshape(Pxx_pre.shape[0],-1)

    Pxx_2D = TSNE(2).fit_transform(np.concatenate((Pxx_cur,Pxx_pre),axis=0))
    Prob = Prob_pre.reshape(-1,1)
    Prob = Scale_Prob(Prob,scale=(0,1))
    colors = ['cornflowerblue','blueviolet','slateblue','plum']
    plt.style.use("ggplot")
    plt.figure(dpi=dpi,figsize=(5,5))
    for i in range(1,Pxx_2D.shape[0]):
        if Prob[i - 1] > 0.75:
            c = 0
        elif Prob[i - 1] > 0.5:
            c = 1
        elif Prob[i - 1] > 0.25:
            c = 2
        else:
            c = 3
        f1 = plt.scatter(Pxx_2D[i:,0],Pxx_2D[i:,1],c=colors[c],edgecolors=colors[c])
    f2 = plt.scatter(Pxx_2D[0,0],Pxx_2D[0,1],c='red',marker="*",s=100)
    plt.legend([f1,f2],[C,'Yours'],fontsize=10,loc='upper right')
    plt.title("2D Distribution Map",fontsize=20)
    fig = plt.gcf()
    fig.set_size_inches(14 / 3, 13 / 3)
    plt.savefig(save_path,dpi=dpi)
    if plot:
        plt.show()

def BandPSD(C,save,raw=None,fname=None,start=0.5,duration=None,freq_bands=None,baseline=None,AvgPSD_root='./images',cmap='jet',dpi=600,show=False):
    assert not (raw is None) or not (fname is None)
    import numpy as np
    import matplotlib.font_manager as fm
    from PIL import Image,ImageDraw,ImageFont
    from utils import reshape_array

    if (raw is None) and not (fname is None):
        try:
            raw = mne.io.read_raw_eeglab(fname, preload=False)
        except:
            print("Error: Read from %s failed!"%os.path.basename(fname))
            raise
    if (duration is None):
        duration = raw._last_time - 0.5 - start
    else:
        assert duration <= raw._last_time - start
    try:
        events = mne.make_fixed_length_events(raw,duration=duration,start=start)
        picks = mne.pick_types(raw.info, meg=False, eeg=True, stim=True, eog=True, exclude='bads')
        epochs = mne.Epochs(raw, events, event_id=1, tmin=start,tmax=duration , proj=True,
                            picks=picks,baseline=baseline, preload=True)
        if (freq_bands is None):
            freq_bands = [(1, 4, 'Delta'), (4, 8, 'Theta'), (8, 13, 'Alpha'),
                    (13, 30, 'Beta'), (30, 45, 'Gamma')]
        epochs.plot_psd_topomap(ch_type='eeg',bands=freq_bands, normalize=True, cmap=cmap,show=False)
        fig = plt.gcf()
        fig.set_size_inches(75 / 3, 10 / 3)
        # plt.savefig('temp.jpg',dpi=dpi)
        canvas = FigureCanvasAgg(fig)
        canvas.draw() 
        img = np.array(canvas.renderer.buffer_rgba())
        Image.fromarray(np.rot90(img)).convert("RGB").save(fp=save)
    except:
        print("Error: plot PSD topomap failed!")
        raise
    return

    img = np.array(Image.open('temp.jpg'))
    _,W,_=img.shape
    n_bands = len(freq_bands)
    img_reshape = img[:,:int(W/n_bands),:]
    for i in range(1,n_bands):
        img_reshape = np.concatenate((img_reshape,img[:,int(i*W/n_bands):int((i+1)*W/n_bands),:]),axis=0)
    img_reshape = np.concatenate((np.ones((100, img_reshape.shape[1], img_reshape.shape[2]), dtype=np.uint8) * 255, img_reshape), axis=0)
    img_pil = Image.fromarray(img_reshape)
    ttfront = ImageFont.truetype(fm.findfont(fm.FontProperties(family='DejaVu Sans')), 90)  # 字体大小
    draw = ImageDraw.Draw(img_pil)
    draw.text(xy=(50, 20), text="Power spectrum topology diagram", font=ttfront, fill=(0, 15, 15))
    avg_psd = np.array(Image.open(os.path.join('.',AvgPSD_root,'avg_psd_%s.jpg'%C)))
    avg_psd = reshape_array(avg_psd,np.array(img_pil))
    img_cat = Image.fromarray(np.concatenate((np.array(img_pil),avg_psd),axis=1))
    img_cat.save(fp=save)
    if show:
        img_pil.show()
    if os.path.exists('./temp.jpg'):
        os.remove('temp.jpg')

def plot_bands_histogram(save,raw=None,fname=None,txt_path=None, channels:list=None,C='HC',freq_band=(1,44),dpi=200,show=False):
    assert not (raw is None) or not (fname is None)
    assert txt_path
    from utils import CalPxx,AvgSpectrum_channel,PsdProbFromTxt,Relative_PSD
    import numpy as np
    import seaborn,os
    from PIL import Image,ImageFont,ImageDraw
    from utils import MyFormatter,y_update_scale_value
    import matplotlib.font_manager as fm

    if not (raw is None):
        Pxx_cur, freq = CalPxx(raw=raw, freq_band=freq_band)
    else:
        Pxx_cur, freq = CalPxx(fname=fname, freq_band=freq_band)

    if raw.info['nchan'] == 32:
        selected_channels = [raw.ch_names.index(_) for _ in ['Fp1', 'Fp2', 'F3', 'F4', 'C3', 'C4', 'P3', 'P4', 'O1',
                                                             'O2', 'F7', 'F8', 'T7', 'T8', 'P7','P8']]
        Pxx_cur = Pxx_cur.squeeze()[selected_channels,:]
    Pxx_cur = Relative_PSD(Pxx_cur)
    Pxx_cur_avg = AvgSpectrum_channel(Pxx_cur) ## N x channels x bands
    Pxx_pre = PsdProbFromTxt(fname=os.path.join('.',txt_path.replace('.','_16channels.' if raw.info['nchan'] >= 16 else '_8channels.')))
    Pxx_pre = Pxx_pre[:,:-1].reshape(Pxx_pre.shape[0],16 if raw.info['nchan'] >= 16 else 8,freq_band[-1]) ##N x channels x bands
    Pxx_pre_avg = AvgSpectrum_channel(Pxx_pre)

    if (channels is None):
        Pxx = np.concatenate((np.mean(Pxx_pre_avg,axis=1),
                              np.mean(Pxx_cur_avg,axis=1)),axis=0)  ## N+1 x bands
    else:
        Pxx = np.concatenate((np.mean(Pxx_pre_avg[:,channels,:], axis=1),
                              np.mean(Pxx_cur_avg[:,channels,:], axis=1)), axis=0)  ## N+1 x bands
    flag = True
    img = None
    plt.figure(dpi=dpi)
    plt.title("adafa")
    titles = ['Delta','Theta','Alpha','Beta','Gamma']
    ratio = [[190, 370, 210, 770, 2150], [190, 330, 190, 650, 850]]  ##[HC DP]
    C_ratio = ratio[0 if C == 'HC' else 1]
    for band_i in range(Pxx.shape[-1]):
        seaborn.distplot(Pxx[:,band_i],bins=20,color='royalblue')
        plt.plot([Pxx[-1,band_i],Pxx[-1,band_i]],[0,plt.axis()[-1]],color='red',linewidth=3)
        # plt.grid(True)
        plt.gca().yaxis.set_major_formatter(MyFormatter(y_update_scale_value, C_ratio[band_i]))
        plt.legend([C,'Yours'])
        plt.xlabel('Relative Power Spectrum',)
        plt.ylabel('Proportion')

        plt.title(titles[band_i])
        fig = plt.gcf()
        fig.set_size_inches(13 / 3, 10 / 3)
        plt.tight_layout()
        plt.savefig('temp.jpg')
        plt.clf()
        if (flag):
            img = np.array(Image.open('temp.jpg'))
            flag=False
        else:
            img = np.concatenate((img,np.array(Image.open('temp.jpg'))),axis=0)
    img = np.concatenate((np.ones((100,img.shape[1],img.shape[2]),dtype=np.uint8)*255,img),axis=0)
    img_pil = Image.fromarray(img)
    ttfront = ImageFont.truetype(fm.findfont(fm.FontProperties(family='DejaVu Sans')), 60)  # 字体大小
    draw = ImageDraw.Draw(img_pil)
    draw.text(xy=(180, 40), text="Frequency Band Power Histogram", font=ttfront, fill=(0, 15, 15))
    img_pil.save(save)
    if show:
        img_pil.show()
    if os.path.exists('./temp.jpg'):
        os.remove('./temp.jpg')

def Pie_plot(raw,save_path:str,model_type='SimpleCNN',double_decision=False,pArray:list=None,dpi=300,show=False):

    if not (pArray is None):
        assert len(pArray) == 2
        if double_decision:
            values = [1-pArray[0],pArray[0],1-pArray[1],pArray[1]]
        else:
            values = [pArray[0],pArray[1]]
    else:
        pArray = []
        from models import predict
        for model_type in ['SimpleCNN', 'RF']:
            pArray += predict(raw=raw,model_type=model_type).tolist()
        if double_decision:
            values = pArray
        else:
            values = pArray[:2]
    C = 'HC' if values[0] > 0.5 else 'DP'

    plt.figure(dpi=dpi)
    f, axes = plt.subplots(1, 2 if double_decision else 1)
    f.suptitle('AI Diagnosis',fontsize=20)
    labels=['Health','Depression']
    titles = ['Deep Learning','Machine Learning']
    explode=[0.005] * 2
    n_plots = 2 if double_decision else 1
    for i in range(n_plots):
        val = values[2*i:2*i+2]
        if val[0] < 0.01:
            val[0] += 0.01
            val[1] -= 0.01
        elif val[0] > 0.99:
            val[0] -= 0.01
            val[1] += 0.01
        ax = axes[i]
        patches,l_text,p_text = ax.pie(val,explode=explode,autopct='%1.1f%%',radius=0.9,
                                       labeldistance=0.9,colors=['royalblue','coral'])
        for t in l_text:
            t.set_size(8)
        for t in p_text:
            t.set_size(12)
        ax.set_title(titles[i],fontsize=14)
        ax.legend(['Health', 'Depression'],loc='upper right')

    plt.tight_layout()
    plt.subplots_adjust(left=0, bottom=0, right=1, top=0.9,hspace = 0.1, wspace = 0.1)
    fig = plt.gcf()
    fig.set_size_inches(22 / 3, 13 / 3)
    plt.savefig(save_path,dpi=dpi)
    if show:
        plt.show()
    return C


# if __name__ == '__main__':
#     path = 'zhangsan_005_15612345678.set'
#     raw = mne.io.read_raw_eeglab(path)
#     TimeSeries_plot(save_path='TimeSeries.jpg', raw=raw, start_time=0.5, duration=3, dpi=300, show=False)
#     TSNE_plot(C='HC', save_path='TSNE.jpg', raw=raw, fname=None, txt_path='HC.txt', freq_bands=(1,44), dpi=300, plot=False)
#     BandPSD(C='HC',save='PSD_Topomap.jpg', raw=raw, fname=None, start=1, duration=None,baseline=None, cmap='jet',dpi=300,show=False)
#     Pie_plot(raw=raw,save_path='./Pie_double.jpg',double_decision=True,dpi=300)
#     plot_bands_histogram(save='Band_histogram.jpg', raw=raw, fname=None, txt_path='HC.txt',
#                          channels= None, C = 'HC', dpi = 300, show = False)
