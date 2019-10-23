import {createStore} from 'redux';
import {storeActive} from "./Reducers"
const store=createStore(storeActive);

export default store