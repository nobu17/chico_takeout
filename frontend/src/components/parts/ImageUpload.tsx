import * as React from "react";
import { Stack, TextField, Button } from "@mui/material";
import { v4 as uuids4 } from "uuid";
import { getImageUrl } from "../../libs/util/ImageUtil"

export type ImageUploadProps = {
  imageUrl: string;
  onImageUploaded: callbackType;
  imageUploading: boolean;
};

interface callbackType {
  (file: File | null): void;
}

const imageStyle = {
  width: "200px",
  height: "200px",
  backgroundSize: 'contain',
  padding: 0,
  margin: 0,
};

export default function ImageUpload(props: ImageUploadProps) {
  const [uuid] = React.useState(uuids4());
  const message = props.imageUploading ? "アップロード中" : "画像選択";
  const handleButtonClick = () => {
    const up = document.getElementById(uuid);
    up?.click();
  };
  const handleOnChange = (files: FileList | null) => {
    if (files !== null) {
      const file = files[0];
      props.onImageUploaded(file);
    }
  };
  return (
    <>
      <Stack spacing={2}>
        <img
          loading="lazy"
          src={getImageUrl(props.imageUrl)}
          style={imageStyle}
          alt="product"
        />
        <TextField
          value={props.imageUrl}
          InputProps={{
            readOnly: true,
          }}
        ></TextField>
        <Button
          disabled={props.imageUploading}
          onClick={handleButtonClick}
          color="warning"
          variant="contained"
          size="large"
        >
          {message}
          <input
            id={uuid}
            type="file"
            hidden
            accept="image/jpeg, image/jpg, image/png"
            onChange={(e) => handleOnChange(e.target.files)}
          />
        </Button>
      </Stack>
    </>
  );
}
