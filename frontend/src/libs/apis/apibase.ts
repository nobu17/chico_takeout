import axios, { AxiosInstance, AxiosError } from "axios";

export interface ApiResponse<T> {
  data: T;
}

interface ApiErrorImpl {
  error: Error;
  code: number;
  message: string;
  isBadRequest: () => boolean;
}

export class ApiError implements ApiErrorImpl {
  public name: string;
  constructor(public error: Error, public code: number, public message: string) {
    this.name = "ApiError";
  }

  isBadRequest(): boolean {
    return this.code === 400;
  }
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
          reject(convertError(error));
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
          reject(convertError(error));
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
          reject(convertError(error));
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
          reject(convertError(error));
        });
    });
  }
}

const convertError = (error: any): ApiError => {
  if (isAxiosError(error)) {
    const message =
      typeof error.response?.data === "string" ? error.response?.data : "";
    return new ApiError(
      error,
      error.response?.status != null ? error.response?.status : 0,
      message,
    );
  }
  return new ApiError(
    error,
    0,
    "",
  );
};

const isAxiosError = (error: any): error is AxiosError => {
  return !!error.isAxiosError;
};
