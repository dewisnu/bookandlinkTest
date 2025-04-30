import React, { useState, useEffect } from 'react';
import { RefreshCw, Filter } from 'lucide-react';
import Header from './components/Header';
import JobsTable from './components/JobsTable';
import UploadForm from './components/UploadForm';
import Pagination from './components/Pagination';
import { Job } from './types';
import { fetchJobs, fetchJobsByStatus,retryJob } from './api/jobs';

function App() {
  const [jobs, setJobs] = useState<Job[]>([]); // Ensure jobs is always an array
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshing, setRefreshing] = useState(false);
  const [statusFilter, setStatusFilter] = useState<string | null>(null);

  // Pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 5;

  // Initial load and periodic refresh
  useEffect(() => {
    const loadJobs = async () => {
      try {
        setLoading(true);
        const data = statusFilter
            ? await fetchJobsByStatus(statusFilter)
            : await fetchJobs();
        if (data === null) {
          throw new Error('Failed to fetch jobs: data is undefined');
        }
        setJobs(data);
        setError(null);
      } catch (err) {
        setError('Failed to load jobs. Please try again.');
        console.error('Error loading jobs:', err);
      } finally {
        setLoading(false);
      }
    };

    loadJobs();

    // Set up polling
    const interval = setInterval(loadJobs, 5000);

    return () => clearInterval(interval);
  }, [statusFilter]);

  // Get current page's jobs
  const indexOfLastItem = currentPage * itemsPerPage;
  const indexOfFirstItem = indexOfLastItem - itemsPerPage;
  const currentJobs = (jobs || []).slice(indexOfFirstItem, indexOfLastItem); // safeguard for null or undefined

  const handlePageChange = (pageNumber: number) => {
    setCurrentPage(pageNumber);
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      const data = statusFilter
          ? await fetchJobsByStatus(statusFilter)
          : await fetchJobs();

      if (data === null) {
        throw new Error('Failed to fetch jobs: data is undefined');
      }

      setJobs(data);
      setError(null);
    } catch (err) {
      setError('Failed to refresh jobs. Please try again.');
      console.error('Error refreshing jobs:', err);
    } finally {
      setRefreshing(false);
    }
  };

  const handleUploadSuccess = async () => {
    await handleRefresh();
    setCurrentPage(1); // Reset to first page after new upload
  };

  const handleFilterChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setStatusFilter(e.target.value || null);
    setCurrentPage(1); // Reset to first page when filter changes
  };

  return (
      <div className="min-h-screen bg-gray-50">
        <Header />
        <main className="container mx-auto px-4 py-8">
          <div className="mb-8">
            <UploadForm onSuccess={handleUploadSuccess} />
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
              <h2 className="text-xl font-semibold text-gray-800">Job History</h2>
              <div className="flex flex-col sm:flex-row gap-3 w-full md:w-auto">
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                    <Filter size={16} className="text-gray-400" />
                  </div>
                  <select
                      value={statusFilter || ''}
                      onChange={handleFilterChange}
                      className="block w-full pl-10 py-2 pr-3 text-sm border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="">All Status</option>
                    <option value="pending">Pending</option>
                    <option value="processing">Processing</option>
                    <option value="complete">Complete</option>
                    <option value="failed">Failed</option>
                  </select>
                </div>

                <button
                    onClick={handleRefresh}
                    className="flex items-center justify-center gap-2 bg-blue-50 text-blue-600 px-4 py-2 rounded-md hover:bg-blue-100 transition-colors sm:w-auto w-full"
                    disabled={refreshing}
                >
                  <RefreshCw size={16} className={refreshing ? "animate-spin" : ""} />
                  Refresh
                </button>
              </div>
            </div>

            {error && (
                <div className="bg-red-50 text-red-600 p-4 rounded-md mb-4">
                  {error}
                </div>
            )}

            <div className="mb-3 text-sm text-gray-500">
              Showing {jobs.length} {jobs.length === 1 ? 'job' : 'jobs'}
              {statusFilter ? ` with status "${statusFilter}"` : ''}
            </div>

            <JobsTable
                jobs={currentJobs}
                loading={loading}
                startIndex={(currentPage - 1) * itemsPerPage}
                onRetry={retryJob}
            />

            <Pagination
                currentPage={currentPage}
                totalItems={jobs.length}
                itemsPerPage={itemsPerPage}
                onPageChange={handlePageChange}
            />
          </div>
        </main>
      </div>
  );
}

export default App;

