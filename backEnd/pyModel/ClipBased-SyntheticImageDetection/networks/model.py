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

class CLIPWithSelfAttention(nn.Module):
    """简化版CLIP模型，仅集成自注意力机制"""
    
    def __init__(self, feature_extractor, embedding_dim=512, num_classes=2, dropout_rate=0.3):
        super().__init__()
        self.feature_extractor = feature_extractor
        for param in self.feature_extractor.parameters():
            param.requires_grad = False
            
        # 自注意力模块
        self.self_attention = SelfAttention(embedding_dim)
        
        # 新型分类头
        self.classifier = nn.Sequential(
            nn.Linear(embedding_dim, embedding_dim),
            nn.LayerNorm(embedding_dim),
            nn.ReLU(),
            nn.Dropout(dropout_rate),
            nn.Linear(embedding_dim, num_classes)
        )

    def forward(self, x):
        # 特征提取
        with torch.no_grad():
            features = self.feature_extractor(x)  # 假设输出形状 [batch, channels, h, w]
            
        # 转换为适合注意力的格式 [batch, seq_len, embedding_dim]
        batch_size, C, H, W = features.shape
        features = features.view(batch_size, C, -1).permute(2, 0, 1)  # [H*W, batch, C]
        
        # 应用自注意力
        attn_features = self.self_attention(features)
        
        # 池化操作
        pooled = attn_features.mean(dim=0)  # [batch, C]
        
        # 分类
        return self.classifier(pooled)