import React from 'react';
import { Image } from 'lucide-react';

const Header: React.FC = () => {
  return (
    <header className="bg-white shadow sticky top-0 z-10">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="p-2 bg-blue-50 rounded-md">
              <Image size={24} className="text-blue-600" />
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-800">Image Processing Dashboard</h1>
              <p className="text-sm text-gray-500">Monitor and manage image compression jobs</p>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;