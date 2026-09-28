import { composeWithDevTools } from '@redux-devtools/extension';
import { legacy_createStore as createStore } from 'redux';

import allReducers from './reducers';

const store = createStore(
  allReducers,
  process.env.NODE_ENV === 'development' ? composeWithDevTools() : undefined
);

export default store;
