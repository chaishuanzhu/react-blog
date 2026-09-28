import type { UnknownAction } from 'redux';

import { SET_NAV_SHOW } from '../constant';

const initState = true;

export default function addReducer(preState = initState, action: UnknownAction) {
  switch (action.type) {
    case SET_NAV_SHOW:
      return action.data as boolean;
    default:
      return preState;
  }
}
