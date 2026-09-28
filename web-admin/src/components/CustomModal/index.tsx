import { Input, Modal } from '@arco-design/web-react';
import React from 'react';

import ModalTitle from '../ModalTitle';
import s from './index.scss';

export interface ModalField {
  text: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}

interface Props {
  isEdit: boolean;
  isModalOpen: boolean;
  name: string;
  modalOk: () => void;
  modalCancel: () => void;
  confirmLoading?: boolean;
  fields?: ModalField[];
  addText?: string;
  updateText?: string;
  children?: React.ReactNode;
}

const CustomModal: React.FC<Props> = ({
  isEdit,
  isModalOpen,
  name,
  modalOk,
  modalCancel,
  confirmLoading,
  fields = [],
  addText = '添加',
  updateText = '更新',
  children
}) => (
  <Modal
    title={
      <ModalTitle isEdit={isEdit} name={name} addText={addText} updateText={updateText} />
    }
    visible={isModalOpen}
    onOk={modalOk}
    onCancel={modalCancel}
    confirmLoading={confirmLoading}
    unmountOnExit
  >
    {children ??
      fields.map(({ text, value, onChange, placeholder }) => (
        <Input
          size='large'
          key={text}
          addBefore={text}
          placeholder={placeholder}
          value={value}
          onChange={onChange}
          className={s.modalInput}
        />
      ))}
  </Modal>
);

export default CustomModal;
