import '@arco-design/web-react/es/_util/react-19-adapter';
import '@arco-design/web-react/dist/css/arco.css';

import { createRoot } from 'react-dom/client';
import { BrowserRouter } from 'react-router';

import App from './App';
import ErrorBoundary from './components/ErrorBoundary';
import { basePath } from './utils/constant';

createRoot(document.getElementById('root')!).render(
  <ErrorBoundary>
    <BrowserRouter basename={basePath}>
      <App />
    </BrowserRouter>
  </ErrorBoundary>
);
