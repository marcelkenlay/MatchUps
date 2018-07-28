import {createStore, applyMiddleware} from "redux";
import thunk from "redux-thunk";
import rootReducer from "./reducers"
import promiseMiddleware from 'redux-promise';
import logger from "redux-logger"


const initialState = {};

const middleware = [thunk, promiseMiddleware,logger];

//createStore(reducer, [preloadedState], [enhancer])
const store = createStore(rootReducer, initialState, applyMiddleware(...middleware));

export default store;