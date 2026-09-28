import { combineReducers } from 'redux';

import avatar from './avatar';
import email from './email';
import link from './link';
import mode from './mode';
import name from './name';
import navShow from './navShow';

export default combineReducers({
  navShow,
  avatar,
  email,
  link,
  name,
  mode
});
