class SelfAttention(nn.Module):
    def __init__(self, in_dim):
        super(SelfAttention, self).__init__()
        self.query = nn.Linear(in_dim, in_dim)
        self.key = nn.Linear(in_dim, in_dim)
        self.value = nn.Linear(in_dim, in_dim)
        self.scale = torch.sqrt(torch.tensor(in_dim, dtype=torch.float32))
        self.gamma = nn.Parameter(torch.zeros(1)) 

class EnhancedClassifier(nn.Module):
    def __init__(self, in_features, hidden_dim=512, num_classes=2, dropout_rate=0.3):
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

class EnhancedCLIPModel(nn.Module):
    def __init__(self, feature_extractor, embedding_dim=512, hidden_dim=512, num_classes=2, dropout_rate=0.3):
        self.feature_extractor = feature_extractor
        # 冻结特征提取器
        for param in self.feature_extractor.parameters():
            param.requires_grad = False 