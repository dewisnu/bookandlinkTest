import React from 'react';
import { CheckCircle, Clock, Loader, AlertTriangle } from 'lucide-react';

type Status = 'pending' | 'processing' | 'complete' | 'failed';

interface StatusBadgeProps {
  status: Status;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ status }) => {
  const getStatusConfig = () => {
    switch (status) {
      case 'pending':
        return {
          icon: <Clock size={12} />,
          classes: 'bg-gray-100 text-gray-800',
          label: 'Pending'
        };
      case 'processing':
        return {
          icon: <Loader size={12} className="animate-spin" />,
          classes: 'bg-yellow-100 text-yellow-800',
          label: 'Processing'
        };
      case 'complete':
        return {
          icon: <CheckCircle size={12} />,
          classes: 'bg-green-100 text-green-800',
          label: 'Complete'
        };
      case 'failed':
        return {
          icon: <AlertTriangle size={12} />,
          classes: 'bg-red-100 text-red-800',
          label: 'Failed'
        };
      default:
        return {
          icon: <Clock size={12} />,
          classes: 'bg-gray-100 text-gray-800',
          label: status.charAt(0).toUpperCase() + status.slice(1)
        };
    }
  };

  const { icon, classes, label } = getStatusConfig();

  const baseClasses = 'px-2 py-1 rounded-full text-xs font-medium inline-flex items-center gap-1';
  
  return (
    <span className={`${baseClasses} ${classes}`}>
      {icon}
      {label}
    </span>
  );
};

export default StatusBadge;