import React from 'react'

//import {createRoot} from 'react-dom/client'

import ReactDOM from "react-dom";


import './style.css'
import App from './App'
import {createRoot} from "react-dom/client"

const container = document.getElementById('root')


const root = createRoot(container)

root.render(
   // <React.StrictMode>
        <App/>
    //</React.StrictMode>
)