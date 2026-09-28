import type { ButtonProps } from '@arco-design/web-react';
import { Button, Message } from '@arco-design/web-react';
import { IconUpload } from '@arco-design/web-react/icon';
import React, { useRef, useState } from 'react';

import { uploadImage } from '@/utils/api';
import { showError } from '@/utils/feedback';

interface Props {
  onUploaded: (url: string, file: File) => void;
  text?: string;
  size?: ButtonProps['size'];
  style?: React.CSSProperties;
}

const UploadButton: React.FC<Props> = ({ onUploaded, text = '上传图片', size, style }) => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);

  const handleChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    e.target.value = '';
    if (!file) return;
    setUploading(true);
    try {
      const url = await uploadImage(file);
      onUploaded(url, file);
      Message.success('上传成功！');
    } catch (err) {
      showError(err);
    } finally {
      setUploading(false);
    }
  };

  return (
    <>
      <input
        ref={inputRef}
        type='file'
        accept='image/jpeg,image/png,image/gif,image/webp,image/avif'
        style={{ display: 'none' }}
        onChange={handleChange}
      />
      <Button
        type='primary'
        size={size}
        style={style}
        icon={<IconUpload />}
        loading={uploading}
        onClick={() => inputRef.current?.click()}
      >
        {text}
      </Button>
    </>
  );
};

export default UploadButton;
