import type { UnknownAction } from 'redux';

import { SET_EMAIL } from '../constant';

const initState = '';

export default function addReducer(preState = initState, action: UnknownAction) {
  switch (action.type) {
    case SET_EMAIL:
      return action.data as string;
    default:
      return preState;
  }
}
