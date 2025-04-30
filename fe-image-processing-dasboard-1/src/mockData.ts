import { Job } from './types';

// Generate a timestamp between 1 and 24 hours ago
const getRandomPastTime = () => {
  const now = new Date();
  const hoursAgo = Math.floor(Math.random() * 24) + 1;
  now.setHours(now.getHours() - hoursAgo);
  return now.toISOString();
};

// Generate a random file size between min and max KB
const getRandomFileSize = (min: number, max: number) => {
  return Math.floor(Math.random() * (max - min + 1) + min) * 1024; // Convert to bytes
};

// Mock data for image processing jobs
export const mockJobs: Job[] = [
  {
    id: 1,
    filename: 'beach_sunrise.jpg',
    original_size: getRandomFileSize(2000, 5000),
    compressed_size: getRandomFileSize(500, 1500),
    compressed_url: '/images/beach_sunrise_compressed.jpg',
    status: 'complete',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  },
  {
    id: 2,
    filename: 'mountain_view.png',
    original_size: getRandomFileSize(3000, 8000),
    compressed_size: getRandomFileSize(1000, 2000),
    compressed_url: '/images/mountain_view_compressed.png',
    status: 'complete',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  },
  {
    id: 3,
    filename: 'city_skyline.jpg',
    original_size: getRandomFileSize(4000, 6000),
    status: 'processing',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  },
  {
    id: 4,
    filename: 'product_photo.png',
    original_size: getRandomFileSize(1500, 3000),
    status: 'pending',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  },
  {
    id: 5,
    filename: 'family_portrait.jpg',
    original_size: getRandomFileSize(7000, 10000),
    status: 'failed',
    error_message: 'File format not supported',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  },
  {
    id: 6,
    filename: 'sunset_beach.jpg',
    original_size: getRandomFileSize(3000, 5000),
    compressed_size: getRandomFileSize(500, 1500),
    compressed_url: '/images/sunset_beach_compressed.jpg',
    status: 'complete',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  },
  {
    id: 7,
    filename: 'hotel_lobby.png',
    original_size: getRandomFileSize(8000, 12000),
    status: 'processing',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  },
  {
    id: 8,
    filename: 'restaurant_menu.jpg',
    original_size: getRandomFileSize(1000, 2000),
    compressed_size: getRandomFileSize(300, 800),
    compressed_url: '/images/restaurant_menu_compressed.jpg',
    status: 'complete',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  },
  {
    id: 9,
    filename: 'error_large.tiff',
    original_size: getRandomFileSize(20000, 30000),
    status: 'failed',
    error_message: 'File size exceeds maximum limit',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  },
  {
    id: 10,
    filename: 'hotel_room.jpg',
    original_size: getRandomFileSize(5000, 8000),
    compressed_size: getRandomFileSize(1000, 3000),
    compressed_url: '/images/hotel_room_compressed.jpg',
    status: 'complete',
    created_at: getRandomPastTime(),
    updated_at: getRandomPastTime()
  }
];