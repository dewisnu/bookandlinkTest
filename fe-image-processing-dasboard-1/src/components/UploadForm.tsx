import React, { useState, useRef } from 'react';
import { Upload, X, Image, FileText } from 'lucide-react';
import { uploadImages } from '../api/jobs';

interface UploadFormProps {
  onSuccess: () => void;
}

const UploadForm: React.FC<UploadFormProps> = ({ onSuccess }) => {
  const [files, setFiles] = useState<File[]>([]);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [dragActive, setDragActive] = useState(false);
  const [uploadProgress, setUploadProgress] = useState<number>(0);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      const filesArray = Array.from(event.target.files);
      const imageFiles = filesArray.filter(file => file.type.startsWith('image/'));
      
      if (imageFiles.length === 0) {
        setError('Please select image files only (jpg, png, gif, webp)');
        return;
      }
      
      setFiles((prevFiles) => [...prevFiles, ...imageFiles]);
      setError(null);
    }
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(true);
  };

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    
    if (e.dataTransfer.files) {
      const filesArray = Array.from(e.dataTransfer.files);
      const imageFiles = filesArray.filter(file => file.type.startsWith('image/'));
      
      if (imageFiles.length === 0) {
        setError('Please drop image files only (jpg, png)');
        return;
      }
      
      setFiles((prevFiles) => [...prevFiles, ...imageFiles]);
      setError(null);
    }
  };

  const handleRemoveFile = (index: number) => {
    setFiles((prevFiles) => prevFiles.filter((_, i) => i !== index));
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getFileIcon = (file: File) => {
    if (file.type.startsWith('image/')) {
      return <Image size={16} className="text-blue-500" />;
    }
    return <FileText size={16} className="text-gray-500" />;
  };

  const handleUpload = async () => {
    if (files.length === 0) {
      setError('Please select at least one image to upload');
      return;
    }

    setError(null);
    setUploading(true);
    setUploadProgress(0);

    try {
      await uploadImages(files);
      setFiles([]);
      onSuccess();
    } catch (err) {
      setError('Failed to upload images. Please try again.');
      console.error('Upload error:', err);
    } finally {
      setUploading(false);
      setUploadProgress(0);
    }
  };

  const triggerFileInput = () => {
    fileInputRef.current?.click();
  };

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-xl font-semibold text-gray-800 mb-4">Upload Images</h2>
      
      {error && (
        <div className="bg-red-50 text-red-600 p-4 rounded-md mb-4 flex items-center">
          <X size={16} className="mr-2 flex-shrink-0" />
          <span>{error}</span>
        </div>
      )}

      <div 
        className={`border-2 border-dashed rounded-lg p-8 text-center transition-all duration-200 ${
          dragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-blue-400 hover:bg-gray-50'
        }`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        <input
          type="file"
          multiple
          accept="image/*"
          className="hidden"
          onChange={handleFileChange}
          ref={fileInputRef}
        />

        <Upload size={36} className={`mx-auto mb-3 transition-colors duration-200 ${dragActive ? 'text-blue-500' : 'text-gray-400'}`} />
        <p className="text-gray-700 mb-2">Drag and drop images here or</p>
        <button
          type="button"
          onClick={triggerFileInput}
          className="text-blue-600 font-medium hover:text-blue-800 transition-colors duration-200"
        >
          Browse files
        </button>
        <p className="text-sm text-gray-500 mt-2">
          Supported formats: JPG, PNG
        </p>
      </div>

      {files.length > 0 && (
        <div className="mt-4">
          <h3 className="font-medium text-gray-700 mb-2">Selected files ({files.length})</h3>
          <ul className="space-y-2 max-h-40 overflow-y-auto rounded-md border border-gray-200">
            {files.map((file, index) => (
              <li 
                key={`${file.name}-${index}`} 
                className="flex items-center justify-between bg-gray-50 hover:bg-gray-100 p-3 border-b border-gray-200 last:border-b-0 transition-colors duration-200"
              >
                <div className="flex items-center space-x-3 truncate max-w-[calc(100%-2.5rem)]">
                  {getFileIcon(file)}
                  <div className="truncate">
                    <p className="text-sm truncate font-medium">{file.name}</p>
                    <p className="text-xs text-gray-500">{formatFileSize(file.size)}</p>
                  </div>
                </div>
                <button 
                  onClick={() => handleRemoveFile(index)}
                  className="text-gray-400 hover:text-red-500 transition-colors duration-200 rounded-full p-1 hover:bg-red-50"
                >
                  <X size={16} />
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}

      {uploading && (
        <div className="mt-4">
          <div className="flex justify-between text-sm text-gray-600 mb-1">
            <span>Uploading...</span>
            <span>{uploadProgress}%</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className="bg-blue-600 h-2 rounded-full transition-all duration-300 ease-out"
              style={{ width: `${uploadProgress}%` }}
            ></div>
          </div>
        </div>
      )}

      <div className="mt-6">
        <button
          onClick={handleUpload}
          disabled={uploading || files.length === 0}
          className={`flex items-center justify-center gap-2 w-full py-3 px-4 rounded-md font-medium transition-all duration-200 ${
            uploading || files.length === 0
              ? 'bg-gray-200 text-gray-500 cursor-not-allowed'
              : 'bg-blue-600 text-white hover:bg-blue-700 shadow hover:shadow-md'
          }`}
        >
          {uploading ? (
            <>
              <div className="animate-spin h-4 w-4 border-2 border-white border-t-transparent rounded-full"></div>
              Uploading...
            </>
          ) : (
            <>
              <Upload size={16} />
              Upload {files.length > 0 ? `${files.length} file${files.length > 1 ? 's' : ''}` : 'images'}
            </>
          )}
        </button>
      </div>
    </div>
  );
};

export default UploadForm;