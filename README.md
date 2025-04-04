# misterblast-storage

1. Upload File
    - Upload if in Server and GCS Empty
    - Upload if in Server empty but in GCS exists
    - Upload if in Server exists but in GCS empty
    - Upload if Key not incuded in payload
    - Upload if API_KEY not included in payload
    - Upload if Image Size > 5 Mib
    - Upload if in Server and GCS exists
    - Upload if Key incorrect
    - Upload Time > 30 seconds
2. Get File
    - Get if in Server exists and GCS exists
    - Get if in Server empty but GCS exists
    - Get if in server exists but GCS empty  
    - Get if Key not incuded in payload
    - Get if API_KEY not included in payload
    - Get if Key incorrect
    - Get Time > 30 seconds
3. Delete File
