import React, { useState } from 'react';
import {AlertCircle, RefreshCw,Download } from 'lucide-react';
import { Job } from '../types';
import StatusBadge from './StatusBadge';

interface JobsTableProps {
  jobs: Job[];
  loading: boolean;
  onRetry?: (id: number) => void;
  startIndex: number;
}

const JobsTable: React.FC<JobsTableProps> = ({ jobs, loading, onRetry, startIndex }) => {
  const [expandedRow, setExpandedRow] = useState<number | null>(null);

  const downloadFile = (url:string, filename:any) => {
    fetch(url, {
      method: "GET",
      headers: {}
    })
        .then(response => response.arrayBuffer())
        .then(buffer => {
          const blob = new Blob([buffer]);
          const downloadUrl = window.URL.createObjectURL(blob);
          const link = document.createElement("a");
          link.href = downloadUrl;
          link.setAttribute("download", filename); // Dynamic filename
          document.body.appendChild(link);
          link.click();
          link.remove(); // Clean up
        })
        .catch(err => {
          console.error("Download error:", err);
        });
  };

  const formatSize = (bytes?: number) => {
    if (bytes === undefined || bytes === null) return '-';
    if (bytes < 1024) return `${bytes} B`;
    const kb = bytes / 1024;
    if (kb < 1024) return `${kb.toFixed(1)} KB`;
    const mb = kb / 1024;
    return `${mb.toFixed(1)} MB`;
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const calculateSavings = (original?: number, compressed?: number) => {
    if (!original || !compressed) return '-';
    if (original === 0) return '0%';
    const savings = ((original - compressed) / original) * 100;
    return `${savings.toFixed(1)}%`;
  };

  const toggleRowExpand = (id: number) => {
    setExpandedRow(expandedRow === id ? null : id);
  };

  if (loading) {
    return (
      <div className="flex justify-center py-8">
        <div className="animate-pulse flex flex-col w-full">
          <div className="h-8 bg-gray-200 rounded w-full mb-4"></div>
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-16 bg-gray-100 rounded w-full mb-2"></div>
          ))}
        </div>
      </div>
    );
  }

  if (jobs.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        <p>No image jobs found. Upload some images to get started.</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              #
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Filename
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Created At
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Original Size
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Compressed Size
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Savings
            </th>
            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {jobs.map((job, index) => (
            <React.Fragment key={job.id}>
              <tr 
                className={`hover:bg-gray-50 transition-colors cursor-pointer ${expandedRow === job.id ? 'bg-blue-50' : ''}`}
                onClick={() => toggleRowExpand(job.id)}
              >
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-500">
                  {startIndex + index + 1}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-800">
                  {job.filename}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  {formatDate(job.created_at)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <StatusBadge status={job.status} />
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  {formatSize(job.original_size)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  {formatSize(job.compressed_size)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                  {calculateSavings(job.original_size, job.compressed_size)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm" onClick={(e) => e.stopPropagation()}>
                  {job.compressed_file_name != null && job.status === 'completed' ? (
                      <button
                          onClick={() =>
                              downloadFile(
                                  `http://localhost:8080/images-compressed/${job.compressed_file_name}`,
                                  job.compressed_file_name
                              )
                          }
                          className="inline-flex items-center gap-1 text-blue-600 hover:text-blue-800"
                      >
                        <Download size={14} />
                        Download
                      </button>

                  ) : job.status === 'failed' ? (
                    <button 
                      className="inline-flex items-center gap-1 text-blue-600 hover:text-blue-800"
                      onClick={() => onRetry && onRetry(job.id)}
                    >
                      <RefreshCw size={14} />
                      Retry
                    </button>
                  ) : (
                    <span className="text-gray-400">-</span>
                  )}
                </td>
              </tr>
              {expandedRow === job.id && (
                <tr className="bg-blue-50">
                  <td colSpan={8} className="px-6 py-4">
                    <div className="text-sm text-gray-700">
                      <h4 className="font-medium mb-2">Job Details</h4>
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <p className="text-xs text-gray-500 mb-1">Job ID</p>
                          <p>{job.id}</p>
                        </div>
                        <div>
                          <p className="text-xs text-gray-500 mb-1">Last Updated</p>
                          <p>{formatDate(job.updated_at)}</p>
                        </div>
                      </div>
                      {job.error_message && (
                        <div className="mt-3 p-3 bg-red-50 rounded border border-red-100 flex items-start gap-2">
                          <AlertCircle size={16} className="text-red-500 mt-0.5" />
                          <span className="text-red-600">{job.error_message}</span>
                        </div>
                      )}
                      {job.status === 'completed' && job.original_size && job.compressed_size && (
                        <div className="mt-3">
                          <p className="text-xs text-gray-500 mb-2">Compression Result</p>
                          <div className="w-full bg-gray-200 rounded-full h-2.5">
                            <div 
                              className="bg-green-600 h-2.5 rounded-full" 
                              style={{ width: `${(job.compressed_size / job.original_size) * 100}%` }}
                            ></div>
                          </div>
                          <div className="flex justify-between text-xs mt-1">
                            <span className="text-green-600">{formatSize(job.compressed_size)}</span>
                            <span className="text-gray-500">{formatSize(job.original_size)}</span>
                          </div>
                        </div>
                      )}
                    </div>
                  </td>
                </tr>
              )}
            </React.Fragment>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default JobsTable;