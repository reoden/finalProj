'''                                        
Copyright 2024 Image Processing Research Group of University Federico
II of Naples ('GRIP-UNINA'). All rights reserved.
                        
Licensed under the Apache License, Version 2.0 (the "License");       
you may not use this file except in compliance with the License. 
You may obtain a copy of the License at                    
                                           
    http://www.apache.org/licenses/LICENSE-2.0
                                                      
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,    
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.                         
See the License for the specific language governing permissions and
limitations under the License.
''' 

import torch
import torch.nn as nn
import torch.nn.functional as F
from typing import Optional

class SelfAttention(nn.Module):
    """自注意力机制层，用于捕获特征间的长距离依赖关系"""
    
    def __init__(self, in_dim):
        super(SelfAttention, self).__init__()
        self.query = nn.Linear(in_dim, in_dim)
        self.key = nn.Linear(in_dim, in_dim)
        self.value = nn.Linear(in_dim, in_dim)
        self.scale = torch.sqrt(torch.tensor(in_dim, dtype=torch.float32))
        self.gamma = nn.Parameter(torch.zeros(1))
        
    def forward(self, x):
        # x: [batch_size, length, in_dim]
        query = self.query(x)
        key = self.key(x)
        value = self.value(x)
        
        attention = torch.matmul(query, key.transpose(-2, -1)) / self.scale
        attention = F.softmax(attention, dim=-1)
        
        out = torch.matmul(attention, value)
        return self.gamma * out + x

class EnhancedClassifier(nn.Module):
    """增强的分类头，使用多层感知机、注意力机制和Dropout来提高性能"""
    
    def __init__(self, in_features, hidden_dim=512, num_classes=2, dropout_rate=0.3):
        super(EnhancedClassifier, self).__init__()
        self.mlp = nn.Sequential(
            nn.Linear(in_features, hidden_dim),
            nn.LayerNorm(hidden_dim),
            nn.ReLU(),
            nn.Dropout(dropout_rate),
            nn.Linear(hidden_dim, hidden_dim // 2),
            nn.LayerNorm(hidden_dim // 2),
            nn.ReLU(),
            nn.Dropout(dropout_rate)
        )
        self.attention = SelfAttention(hidden_dim // 2)
        self.classifier = nn.Linear(hidden_dim // 2, num_classes)
        
    def forward(self, x):
        # 如果x是特征图(2D或更高维度)，先将其展平
        if len(x.shape) > 2:
            batch_size = x.shape[0]
            x = x.reshape(batch_size, -1)
            
        # 通过MLP网络
        x = self.mlp(x)
        
        # 应用自注意力(需要调整形状)
        if len(x.shape) == 2:
            x = x.unsqueeze(1)  # [batch_size, 1, hidden_dim//2]
            x = self.attention(x)
            x = x.squeeze(1)    # [batch_size, hidden_dim//2]
        
        # 最终分类
        return self.classifier(x)

class EnhancedCLIPModel(nn.Module):
    """增强的CLIP模型，结合了特征提取器和改进的分类头"""
    
    def __init__(self, feature_extractor, embedding_dim=512, hidden_dim=512, num_classes=2, dropout_rate=0.3):
        super(EnhancedCLIPModel, self).__init__()
        self.feature_extractor = feature_extractor
        # 冻结特征提取器
        for param in self.feature_extractor.parameters():
            param.requires_grad = False
            
        self.classifier = EnhancedClassifier(
            embedding_dim, 
            hidden_dim=hidden_dim,
            num_classes=num_classes,
            dropout_rate=dropout_rate
        )
        
    def forward(self, x):
        # 提取特征
        with torch.no_grad():
            features = self.feature_extractor(x)
        
        # 分类
        return self.classifier(features)
    
    def unfreeze_last_layers(self):
        """解冻特征提取器的最后几层以进行微调"""
        # 这是一个示例，实际解冻的层应该根据特定的特征提取器结构定制
        for name, param in self.feature_extractor.named_parameters():
            if "layer4" in name or "layer3" in name:
                param.requires_grad = True 