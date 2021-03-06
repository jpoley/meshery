import fetch from 'isomorphic-unfetch'

const dataFetch = (url, options = {}, successFn, errorFn) => {
  fetch(url, options)
    .then(res => {
      if (res.status === 401 || res.redirected){
        if (window.location.host.endsWith('3000')){
          window.location = "/login"; // for local dev thru node server
        } else {
          window.location.reload(); // for use with Go server
        }
      }
      let result;
      if (res.ok) {
        // console.log(`res type: ${res.type}`);
        try {
          result = res.json();
        } catch(e){
          result = res.text();
        }
        return result;
      } else {
        res.text().then(errorFn);
      }

    }).then(successFn)
    .catch(errorFn);
}

export default dataFetch;