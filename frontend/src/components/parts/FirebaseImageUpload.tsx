import React, { useState } from "react";
import { Alert, Typography } from "@mui/material";
import ImageUpload from "./ImageUpload";
import { uploadImageWithCompress } from "../../libs/firebase/StorageApi";
import { v4 as uuids4 } from "uuid";
import LoadingSpinner from "./LoadingSpinner";

export type FirebaseImageUploadProps = {
  baseUrl: string;
  imageUrl: string;
  onImageUploaded: callbackType;
};

interface callbackType {
  (fileUrl: string): void;
}

export const FirebaseImageUpload = (props: FirebaseImageUploadProps) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error>();
  const { onImageUploaded, ...rest } = props;

  const handleImageUploaded = async (file: File | null) => {
    if (!file) return;

    try {
      setLoading(true);
      const name = uuids4() + ".jpg";
      const newUrl = await uploadImageWithCompress(props.baseUrl, name, file);
      onImageUploaded(newUrl);
    } catch (err: any) {
      console.error("failed to upload file", err);
      setError(err);
    } finally {
      setLoading(false);
    }
  };

  const errorMessage = (error: Error | undefined) => {
    if (error) {
      console.log("err", error);
      return (
        <Alert variant="filled" severity="error">
          エラーが発生しました。
        </Alert>
      );
    }
    return <></>;
  };

  return (
    <>
      <Typography color="text.primary">商品画像</Typography>
      {errorMessage(error)}
      <ImageUpload
        imageUploading={loading}
        onImageUploaded={handleImageUploaded}
        {...rest}
      ></ImageUpload>
      <LoadingSpinner message="Loading..." isLoading={loading} />
    </>
  );
};
