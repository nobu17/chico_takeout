import axios, { AxiosInstance } from "axios";

export interface ApiResponse<T> {
  data: T;
}

export default class ApiBase {
  private _api: AxiosInstance;
  private _baseUrl: string;
  constructor(baseUrl: string) {
    this._baseUrl = baseUrl;
    this._api = axios.create({
      baseURL: baseUrl,
      headers: {
        "Content-Type": "application/json",
        "X-Requested-With": "XMLHttpRequest",
      },
      responseType: "json",
    });
  }

  getAsync<T>(url: string): Promise<ApiResponse<T>> {
    return new Promise((resolve, reject) => {
      const reqUrl = this._baseUrl + url;
      this._api
        .get(reqUrl)
        .then((r) => {
          const res = {
            data: r.data,
          };
          resolve(res);
        })
        .catch((error) => {
          reject(error);
        });
    });
  }

  postAsync(url: string, param: any): Promise<void> {
    const json = JSON.stringify(param);
    return new Promise((resolve, reject) => {
      const reqUrl = this._baseUrl + url;
      this._api
        .post(reqUrl, json)
        .then((r) => {
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  }

  putAsync(url: string, param: any): Promise<void> {
    const json = JSON.stringify(param);
    return new Promise((resolve, reject) => {
      const reqUrl = this._baseUrl + url;
      this._api
        .put(reqUrl, json)
        .then((r) => {
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  }

  deleteAsync(url: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const reqUrl = this._baseUrl + url;
      this._api
        .delete(reqUrl)
        .then((r) => {
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  }
}
