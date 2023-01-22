import torch,mne
import numpy as np
from sklearn.model_selection import GridSearchCV

def ML_score(model,train_x,train_y,test_x,test_y,param=None,n_class=3):
    assert train_x.shape[0] == train_y.shape[0]
    assert test_x.shape[0] == test_y.shape[0]
    if not (param is None):
        model = GridSearchCV(model,param,cv=5)
    model.fit(train_x,train_y)
    train_score = model.score(train_x,train_y)
    test_score = model.score(test_x,test_y)
    test_prediction = model.predict(test_x)
    recall,precision = [],[]
    for i in range(n_class):
        TP = len([_ for _ in range(len(test_y)) if (test_y[_] == i and test_prediction[_] == i)])
        FN = len([_ for _ in range(len(test_y)) if (test_y[_] == i and test_prediction[_] != i)])
        FP = len([_ for _ in range(len(test_y)) if (test_y[_] != i and test_prediction[_] == i)])
        TN = len([_ for _ in range(len(test_y)) if (test_y[_] != i and test_prediction[_] != i)])
        recall.append(float(TP) / (TP+FN+0.0001))
        precision.append(float(TP) / (TP+FP+0.0001))

    return train_score,test_score,recall,precision,test_prediction

class SimpleCNN_16(torch.nn.Module):
    def __init__(self):
        super(SimpleCNN_16, self).__init__()
        self.conv1 = torch.nn.Sequential(
            torch.nn.Conv2d(in_channels=1, out_channels=16, kernel_size=3, stride=1, padding=1),
            torch.nn.BatchNorm2d(16),
            torch.nn.ReLU(),
            torch.nn.AvgPool2d(kernel_size=2, stride=2, padding=0)
        )
        self.conv2 = torch.nn.Sequential(
            torch.nn.Conv2d(in_channels=16, out_channels=32, kernel_size=3, stride=1, padding=1),
            torch.nn.BatchNorm2d(32),
            torch.nn.ReLU(),
        )

        self.flatten = torch.nn.Sequential(
            # 平滑层
            torch.nn.Flatten(),

        )
        self.fc = torch.nn.Sequential(
            torch.nn.Dropout(0.3),
            torch.nn.Linear(32 * 22 * 8, 256),
            torch.nn.ReLU(),
            torch.nn.Linear(256, 2),
            torch.nn.Softmax(dim=1),
        )

    def forward(self, img):
        feature = self.conv1(img)
        feature = self.conv2(feature)
        feature = self.flatten(feature)
        output = self.fc(feature)
        return output

class SimpleCNN_8(torch.nn.Module):
    def __init__(self):
        super(SimpleCNN_8, self).__init__()
        self.conv1 = torch.nn.Sequential(
            torch.nn.Conv2d(in_channels=1, out_channels=16, kernel_size=3, stride=1, padding=1),
            torch.nn.BatchNorm2d(16),
            torch.nn.ReLU(),
            torch.nn.AvgPool2d(kernel_size=2, stride=2, padding=0)
        )
        self.conv2 = torch.nn.Sequential(
            torch.nn.Conv2d(in_channels=16, out_channels=32, kernel_size=3, stride=1, padding=1),
            torch.nn.BatchNorm2d(32),
            torch.nn.ReLU()
        )

        self.flatten = torch.nn.Sequential(
            torch.nn.Flatten(),

        )
        self.fc = torch.nn.Sequential(
            torch.nn.Dropout(0.3),
            torch.nn.Linear(32 * 22 * 4, 256),
            torch.nn.ReLU(),
            torch.nn.Linear(256, 2),
            torch.nn.Softmax(dim=1),
        )

    def forward(self, img):
        feature = self.conv1(img)
        feature = self.conv2(feature)
        feature = self.flatten(feature)
        output = self.fc(feature)
        return output

class MUCHf_Net_16(torch.nn.Module):
    def __init__(self):
        super(MUCHf_Net_16, self).__init__()
        self.conv1 = torch.nn.Sequential(

            torch.nn.Conv2d(in_channels=1, out_channels=16, kernel_size=(5,1), stride=1,padding=(2,0),bias=True),
            # torch.nn.ReLU(),
            # torch.nn.Conv2d(in_channels=16, out_channels=16, kernel_size=(3, 1), stride=1, padding=(1, 0), bias=True),
            torch.nn.BatchNorm2d(16),
            torch.nn.ReLU(),
            torch.nn.Conv2d(in_channels=16,out_channels=32,kernel_size=(1,16),stride=1,padding=0,bias=True),
            torch.nn.BatchNorm2d(32),
            torch.nn.ReLU(),
            torch.nn.AvgPool2d(kernel_size=(2,1), stride=2, padding=0)
        )
        self.conv2 = torch.nn.Sequential(
            torch.nn.Conv2d(in_channels=1, out_channels=16, kernel_size=(3,3), padding=1, stride=1,),
            torch.nn.BatchNorm2d(16),
            torch.nn.ReLU(),
            torch.nn.AvgPool2d(kernel_size=2, stride=2, padding=0),

        )
        self.conv3 = torch.nn.Sequential(
            torch.nn.Conv2d(in_channels=16,out_channels=16,kernel_size=3,stride=1,padding=1,bias=False),
            torch.nn.ReLU(),
            # torch.nn.AvgPool2d(kernel_size=2, stride=2, padding=0),
        )

        self.flatten = torch.nn.Sequential(
            torch.nn.Flatten(),
        )

        self.fc = torch.nn.Sequential(
            torch.nn.Dropout(0.4),
            torch.nn.Linear(16 * 11 * 16, 256),
            torch.nn.Tanh(),
            torch.nn.Linear(256, 2),
            torch.nn.Softmax(dim=1),
        )

    def forward(self, img):
        feature = self.conv1(img)
        feature = torch.transpose(feature,1,3)
        feature = self.conv3(self.conv2(feature))
        feature = self.flatten(feature)
        output = self.fc(feature)
        return output


class MUCHf_Net_8(torch.nn.Module):
    def __init__(self):
        super(MUCHf_Net_8, self).__init__()
        self.conv1 = torch.nn.Sequential(

            torch.nn.Conv2d(in_channels=1, out_channels=32, kernel_size=(3,1), stride=1,padding=(1,0),bias=True),
            torch.nn.ReLU(),
            torch.nn.Conv2d(in_channels=32, out_channels=32, kernel_size=(3, 1), stride=1, padding=(1, 0), bias=True),
            torch.nn.BatchNorm2d(32),
            torch.nn.ReLU(),
            torch.nn.Conv2d(in_channels=32,out_channels=32,kernel_size=(1,8),stride=1,padding=0,bias=True),
            torch.nn.BatchNorm2d(32),
            torch.nn.ReLU(),
            torch.nn.AvgPool2d(kernel_size=(2,1), stride=2, padding=0)
        )
        self.conv2 = torch.nn.Sequential(
            torch.nn.Conv2d(in_channels=1, out_channels=16, kernel_size=(3,3), padding=1, stride=1,),
            torch.nn.BatchNorm2d(16),
            torch.nn.ReLU(),
            torch.nn.AvgPool2d(kernel_size=2, stride=2, padding=0),

        )
        self.conv3 = torch.nn.Sequential(
            torch.nn.Conv2d(in_channels=16,out_channels=16,kernel_size=3,stride=1,padding=1,bias=False),
            torch.nn.ReLU(),
            # torch.nn.AvgPool2d(kernel_size=2, stride=2, padding=0),
        )

        self.flatten = torch.nn.Sequential(
            torch.nn.Flatten(),
        )

        self.fc = torch.nn.Sequential(
            torch.nn.Dropout(0.3),
            torch.nn.Linear(16 * 11 * 16, 256),
            torch.nn.ReLU(),
            torch.nn.Linear(256, 2),
            torch.nn.Softmax(dim=1),
        )

    def forward(self, img):
        feature = self.conv1(img)
        feature = torch.transpose(feature,1,3)
        feature = self.conv3(self.conv2(feature))
        feature = self.flatten(feature)
        output = self.fc(feature)
        return output

class ModelTrainer:
    def __init__(self,chan_num:int,lr,max_epoch,pretained_weight=None):
        # self.model = MUCHf_Net_8() if chan_num == 8 else MUCHf_Net_16()
        self.model = SimpleCNN_16()
        if pretained_weight:
            self.model.load_state_dict(torch.load(pretained_weight))
        else:
            self.__init_params()
        self.lr = lr
        self.max_epoch = max_epoch

    def train(self,x=None,y=None,loader=None):
        from torch.utils.data import TensorDataset,DataLoader
        assert (x and y) or loader

        lr = self.lr
        optimizer = torch.optim.Adam(self.model.parameters(), lr)
        loss_func = torch.nn.CrossEntropyLoss()
        if not loader:
            assert isinstance(x, torch.Tensor) and len(x.size()) == 4
            dataset = TensorDataset(x,y)
            loader = DataLoader(dataset,batch_size=16,shuffle=True)
        train_loss = []
        self.model.train()
        for epoch in range(self.max_epoch):
            if (epoch % 20 == 0):
                print("epoch: %d" % epoch)
            self.adjust_learning_rate(optimizer, epoch, lr)
            n_step = 0
            run_loss = 0
            for x_t, y_t in loader:
                n_step += 1
                out = self.model(x_t)
                loss = loss_func(out, y_t)
                run_loss += loss.item()
                optimizer.zero_grad()
                loss.backward()
                optimizer.step()
            train_loss.append(run_loss / n_step)
        return train_loss

    def test(self,x=None,y=None,loader=None):
        from torch.utils.data import TensorDataset, DataLoader
        assert (x and y) or loader

        if not loader:
            assert isinstance(x, torch.Tensor) and len(x.size()) == 4
            dataset = TensorDataset(x,y)
            loader = DataLoader(dataset,batch_size=16,shuffle=True)
        acc=0
        with torch.no_grad():
            self.model.eval()
            for x_t, y_t in loader:
                out = self.model(x_t)
                out = torch.argmax(out,dim=1)
                acc += sum(out.detach().numpy() == y_t.numpy())

            acc /= loader.dataset.__len__()
        return acc

    def __init_params(self,method='Kaiming'):
        assert method in ['Kaiming','Xavier']
        if method == 'kaiming':
            for module in self.model.modules():
                if isinstance(module,(torch.nn.Conv2d,torch.nn.Linear)):
                    torch.nn.init.kaiming_normal_(module.weight)
                    torch.nn.init.constant_(module.bias.data, 0.0)
                elif isinstance(module, torch.nn.BatchNorm2d):
                    torch.nn.init.constant_(module.weight, 1)
                    torch.nn.init.constant_(module.bias, 0)
        elif method == 'Xavier':
            for module in self.model.modules():
                if isinstance(module, (torch.nn.Conv2d, torch.nn.Linear)):
                    torch.nn.init.xavier_normal_(module.weight,gain=torch.nn.init.calculate_gain('relu'))
                    torch.nn.init.constant_(module.bias.data, 0.0)
                elif isinstance(module, torch.nn.BatchNorm2d):
                    torch.nn.init.constant_(module.weight,1)
                    torch.nn.init.constant_(module.bias, 0)

    def adjust_learning_rate(self,optimizer, epoch, lr):
        """Sets the learning rate to the initial LR decayed by 10 every 30 epochs"""
        lr = lr * (0.8 ** (epoch // 10))
        for param_group in optimizer.param_groups:
            param_group['lr'] = lr

def predict(raw,model_type='SimpleCNN',Pxx=None,freq_band=(1,44)):
    ## Pxx (channels x frequency) relative_power
    import os,joblib
    if not (Pxx is None):
        assert Pxx.shape[-1] == 44 and Pxx.shape[-2] in [8,16]
        if Pxx.max() < 0:
            from utils import Relative_PSD
            Pxx = Relative_PSD(Pxx)
    else:
        from utils import CalPxx,Relative_PSD
        Pxx, _ = CalPxx(raw,freq_band=freq_band)
        assert Pxx.shape[-1] == 44 and Pxx.shape[-2] in [8, 16, 32]
        if Pxx.shape[-2] == 32:
            selected_channels = [raw.ch_names.index(_) for _ in
                                ['Fp1', 'Fp2', 'F3', 'F4', 'C3', 'C4', 'P3', 'P4', 'O1', 'O2', 'F7', 'F8', 'T7', 'T8', 'P7', 'P8']]
            Pxx = Pxx[selected_channels,:]
        Pxx = Relative_PSD(Pxx)
    if 'SimpleCNN' in model_type:
        if Pxx.shape[-2] == 16:
            model = SimpleCNN_16()
        else:
            model = SimpleCNN_8()
        # assert os.path.exists('./SimpleCNN_%dchannels.pth'%Pxx.shape[-2])
        try:
            model.load_state_dict(torch.load('./SimpleCNN_%dchannels.pth'%Pxx.shape[-2]))
        except Exception as e:
            print(e,'Can not load CNN weight!')
            raise
        while len(Pxx.shape) < 4:
            Pxx = np.expand_dims(Pxx, axis=0)
        x = torch.from_numpy(np.transpose(Pxx, [0, 1, 3, 2])).float()
        pred_score = model(x).detach().numpy().squeeze()
    else:
        try:
            if Pxx.shape[-2] == 16:
                model = joblib.load('RF_16channels.pkl')
            else:
                model = joblib.load('RF_8channels.pkl')
        except Exception as e:
            print(e,' Can not load RF')
            raise
        x = np.transpose(Pxx).reshape(1, -1)
        pred_score = model.predict_proba(x).squeeze()
    return pred_score


def BuildDatasetFromRaw(root1,root2,chan_index,freq_band=(1,44),bs=16,test=False):
    import os
    from torch.utils.data import TensorDataset,DataLoader
    from sklearn.model_selection import train_test_split
    from utils import Relative_PSD

    if isinstance(chan_index,str):
        if chan_index == 'all':
            chan_index = list(range(16))
        elif chan_index == '8 channels':
            chan_index = [0,1,2,3,4,5,10,11]
        else:
            chan_index = list(range(16))
    else:
        assert isinstance(chan_index,list) and (len(chan_index) == 8 or len(chan_index) == 16)

    Pxx = []
    target = []
    for root in [root1,root2]:
        length = len(Pxx)
        for file in os.listdir(root):
            if (not file.endswith('.set')):
                continue
            raw = mne.io.read_raw_eeglab(os.path.join(root, file),preload=False)
            Pxx_cur, _ = mne.time_frequency.psd_welch(raw, fmin=freq_band[0], fmax=freq_band[1], n_fft=128, n_overlap=64)
            Pxx.append(np.transpose(Relative_PSD(Pxx_cur)))
        target += [0 if root == root1 else 1] * (len(Pxx) - length)
    Pxx = np.array(Pxx,dtype=np.float32)[:,np.newaxis,:,chan_index]
    target = np.array(target,dtype='long').reshape(Pxx.shape[0],)
    assert list(Pxx.shape[-3:]) == [1,freq_band[-1] - freq_band[0] + 1,len(chan_index)]

    if test:
        p1,p2,t1,t2 = train_test_split(Pxx,target,train_size=0.8,shuffle=True)
        p1 = torch.from_numpy(p1).float()
        p2 = torch.from_numpy(p2).float()
        t1 = torch.from_numpy(t1).long()
        t2 = torch.from_numpy(t2).long()
        dataset1 = TensorDataset(p1,t1)
        train_loader = DataLoader(dataset1,batch_size=bs,shuffle=True)
        dataset2 = TensorDataset(p2, t2)
        test_loader = DataLoader(dataset2, batch_size=bs, shuffle=True)
        return train_loader,test_loader,Pxx,target    ##Pxx(N,1,freq,channel)
    else:
        dataset = TensorDataset(Pxx,target)
        loader = DataLoader(dataset,batch_size=bs,shuffle=True)
        return loader,None,Pxx,target
