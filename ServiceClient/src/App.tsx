import * as React from 'react';
import './static/css/App.css';
import './static/css/iconfont.css';
import { BrowserRouter as BrowserRouter,Switch,Route} from 'react-router-dom';
import Login from './page/Login';
import Home from './page/Home';
class App extends React.Component {
  public render() {
    return (
      <BrowserRouter>
        <Switch>
          <Route path="/home" component={Home} exact={true} />  
          <Route path="/" component={Login}  exact={true} />
        </Switch>
      </BrowserRouter>
    );
  }
}

export default App;
