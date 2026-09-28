import type { UnknownAction } from 'redux';

import { SET_MODE } from '../constant';

const initState = 0;

export default function addReducer(preState = initState, action: UnknownAction) {
  switch (action.type) {
    case SET_MODE:
      return action.data as number;
    default:
      return preState;
  }
}
