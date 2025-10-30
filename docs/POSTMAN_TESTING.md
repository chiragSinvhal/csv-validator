# Postman Testing Guide

## 📋 **Quick Setup**

### 1. Import Collection
1. Open Postman
2. Click **Import**
3. Select `CSV-Validator.postman_collection.json`
4. Import `CSV-Validator.postman_environment.json`

### 2. Set Environment
1. Select **CSV Validator Environment** from the environment dropdown
2. Make sure `base_url` is set to `http://localhost:8080`

## 🚀 **Test Sequence**

### **Step 1: Health Check**
- Run **"Health Check"** request
- ✅ Should return `200 OK` with status
- 🎯 Verifies service is running

### **Step 2: Upload CSV**
- Run **"Upload CSV File"** request
- 📎 Attach a CSV file (use `sample-data/sample1.csv`)
- ✅ Should return `200 OK` with job ID
- 🔄 Job ID automatically saved to environment

### **Step 3: Download Processed File**
- Run **"Download Processed File"** request
- 📥 If ready: Downloads CSV with `has_email` column
- ⏳ If processing: Returns `423 Locked`
- 🎯 Uses job ID from previous step

### **Step 4: Test Error Cases**
- Run **"Upload Invalid File"** (upload .txt file)
- Run **"Download with Invalid Job ID"**
- ✅ Should return proper error responses

## 🎯 **Expected Results**

### **Successful Upload Response**
```json
{
  "id": "a225eb00-0907-4273-92ca-5faadeefae5f"
}
```

### **Processed CSV Output**
```csv
name,email,phone,has_email
Chirag,Chirag@email.com,123-456-7890,true
Jane,,987-654-3210,false
Bob,invalid-email,555-0123,false
```

### **Error Response Example**
```json
{
  "error": "Only CSV files are allowed"
}
```

## 🔧 **Tips for Demo**

1. **Pre-load sample files** in your Downloads folder
2. **Run health check first** to show service is ready
3. **Use sample1.csv** - has mixed email/non-email data
4. **Show both success and error cases**
5. **Explain the async processing** (upload → process → download)

## 📱 **Testing Different Scenarios**

| Scenario | File | Expected Result |
|----------|------|----------------|
| Valid CSV with emails | `sample1.csv` | Success with job ID |
| Valid CSV no emails | `sample2.csv` | Success, all `false` flags |
| Invalid file type | `.txt` file | 400 error |
| Empty file | Empty `.csv` | 400 error |
| Large file | 10MB+ CSV | Success (async processing) |

This collection includes **automated tests** that validate responses, making it perfect for demonstrating reliability to technical reviewers!