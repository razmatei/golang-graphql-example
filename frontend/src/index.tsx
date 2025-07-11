// import { scan } from 'react-scan'; // must be imported before React and React DOM
import React from 'react';
import { createRoot } from 'react-dom/client';
import { RouterProvider, createBrowserRouter } from 'react-router';
import * as dayjs from 'dayjs';
import localizedFormat from 'dayjs/plugin/localizedFormat';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
// import i18n
import './i18n';
import './yup-i18n';
// Import font roboto
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

import subRoutes from './routes/router-routes';
import App from './App';

// if (process.env.NODE_ENV !== 'production') {
//   scan({
//     enabled: true,
//   });
// }

// Extend dayjs
dayjs.extend(localizedFormat);
dayjs.extend(utc);
dayjs.extend(timezone);

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: subRoutes,
  },
]);

const container = document.getElementById('root');
const root = createRoot(container!);
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
