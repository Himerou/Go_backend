async function  request(url, option={}){
    const response = await fetch(url, option);
    if (!response.ok) {
        const data = await response.json();
        console.log(data.error);
        throw new Error(data.error||`HTTP error: ${response.status}`);
    }
    return response;
}

