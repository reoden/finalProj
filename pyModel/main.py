import json
import torch
from transformers import BertTokenizerFast, BertForTokenClassification
from fastapi import FastAPI, UploadFile, HTTPException
from pydantic import BaseModel
import jieba
import pickle

app = FastAPI()

# 定义输入数据模型
class SceneRequest(BaseModel):
    text: str

# 定义输出数据模型
class SceneResponse(BaseModel):
    text: str
    scenes: list[str]

def extract_scene(text, model, tokenizer, device):
    """
    提取文本中的景点信息，确保完整提取词语
    Args:
        text: 输入文本
        model: 加载的模型
        tokenizer: 分词器
        device: 设备类型（CPU/GPU）
    Returns:
        list: 提取出的景点列表
    """
    model.eval()
    # 对输入文本进行编码
    encoding = tokenizer(
        text,
        return_offsets_mapping=True,
        padding=True,
        truncation=True,
        return_tensors="pt"
    )
    
    # 获取input_ids和attention_mask
    input_ids = encoding['input_ids'].to(device)
    attention_mask = encoding['attention_mask'].to(device)
    
    with torch.no_grad():
        outputs = model(input_ids=input_ids, attention_mask=attention_mask)
        predictions = torch.argmax(outputs.logits, dim=2)
    
    # 获取词片段、偏移映射和预测标签
    tokens = tokenizer.convert_ids_to_tokens(input_ids[0])
    offset_mapping = encoding['offset_mapping'][0].numpy()
    label_ids = predictions[0].cpu().numpy()
    
    scenes = []
    current_scene = ''
    scene_start = None
    
    for idx, (token, label_id, offset) in enumerate(zip(tokens, label_ids, offset_mapping)):
        # 跳过特殊token
        if token in ['[CLS]', '[SEP]', '[PAD]'] or offset[0] == offset[1]:
            continue
        
        # 判断是否是景点标签（scene对应的label_id为10）
        if label_id == 10:  # scene对应的label_id
            if scene_start is None:
                scene_start = offset[0]
            
            # 检查下一个token
            next_is_scene = False
            if idx + 1 < len(label_ids):
                next_is_scene = label_ids[idx + 1] == 10
            
            if not next_is_scene:
                # 提取完整的景点名称
                scene_text = text[scene_start:offset[1]]
                if scene_text.strip():
                    scenes.append(scene_text.strip())
                scene_start = None
        else:
            scene_start = None
    
    return scenes

# 加载保存的模型
def load_model(model_path):
    device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
    tokenizer = BertTokenizerFast.from_pretrained(model_path)
    model = BertForTokenClassification.from_pretrained(model_path)
    model.to(device)
    return model, tokenizer, device

def process_file(input_file, output_file, model, tokenizer, device):
    """
    处理整个文件中的文本并保存结果
    """
    results = []
    
    with open(input_file, 'r', encoding='utf-8') as f:
        for line in f:
            data = json.loads(line.strip())
            text = data['text']
            scenes = extract_scene(text, model, tokenizer, device)
            results.append({
                'text': text,
                'scenes': scenes
            })
    
    # 保存结果
    with open(output_file, 'w', encoding='utf-8') as f:
        for item in results:
            json.dump(item, f, ensure_ascii=False)
            f.write('\n')

# 加载模型
model_path = './ner_model'  # 替换为你的模型路径
model, tokenizer, device = load_model(model_path)

@app.post("/extract", response_model=SceneResponse)
def extract_scenes(req: SceneRequest):
  try:
    text = req.text
    scenes = extract_scene(text, model, tokenizer, device)
    return SceneResponse(text=req.text, scenes=scenes)
  
    # 处理文件
    process_file(
        input_file='test.json',
        output_file='scenes_results.json',
        model=model,
        tokenizer=tokenizer,
        device=device
    )
  except Exception as e:
    raise HTTPException(status_code = 500, detail=str(e))    

class MGwordReq(BaseModel):
    text: str

class MGwordResp(BaseModel):
    result: int

def stopwordslist(filepath):
    with open(filepath, 'r', encoding='utf-8') as file:
        stopwords = [line.strip() for line in file.readlines()]
    return stopwords

stopwords = stopwordslist("stopword.txt")

def loadMGword():
    keys = {}
    with open('baiduMGword.txt', encoding='utf-8') as fr:
        for line in fr:
            if line.strip():
                keys[line.strip()] = ''
    return keys  
    
keys = loadMGword()

def tagLabel(text):
    for k in keys:
        if k in text:
            return 1
        else:
            pass
    return 0

@app.post("/scanMGword", response_model=MGwordResp)
def scan_word(req: MGwordReq):
    try:
        text = req.text
        ret = tagLabel(text)
        return MGwordResp(result=ret)
    
    except Exception as e:
        raise HTTPException(status_code = 500, detail=str(e))   

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8888)
