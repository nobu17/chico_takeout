type _Config = {
  apiRoot: string;
};

const Config: _Config = {
  apiRoot: process.env.REACT_APP_API_ROOT || "",
};

export default Config;
