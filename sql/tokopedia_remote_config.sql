PGDMP     '    %                w            tokopedia_remote_config     11.5 (Ubuntu 11.5-1.pgdg18.04+1)     11.5 (Ubuntu 11.5-1.pgdg18.04+1)     r           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                       false            s           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                       false            t           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                       false            u           1262    16429    tokopedia_remote_config    DATABASE     �   CREATE DATABASE tokopedia_remote_config WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';
 '   DROP DATABASE tokopedia_remote_config;
             postgres    false            �            1259    24597    users    TABLE     �   CREATE TABLE public.users (
    id integer NOT NULL,
    first_name character varying(255),
    last_name character varying(255),
    email character varying(255),
    password text
);
    DROP TABLE public.users;
       public         postgres    false            �            1259    24595    users_id_seq    SEQUENCE     �   CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 #   DROP SEQUENCE public.users_id_seq;
       public       postgres    false    197            v           0    0    users_id_seq    SEQUENCE OWNED BY     =   ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;
            public       postgres    false    196            �
           2604    24600    users id    DEFAULT     d   ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);
 7   ALTER TABLE public.users ALTER COLUMN id DROP DEFAULT;
       public       postgres    false    197    196    197            o          0    24597    users 
   TABLE DATA               K   COPY public.users (id, first_name, last_name, email, password) FROM stdin;
    public       postgres    false    197   �
       w           0    0    users_id_seq    SEQUENCE SET     :   SELECT pg_catalog.setval('public.users_id_seq', 6, true);
            public       postgres    false    196            �
           2606    24605    users users_pkey 
   CONSTRAINT     N   ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
 :   ALTER TABLE ONLY public.users DROP CONSTRAINT users_pkey;
       public         postgres    false    197            o   �   x�m�;�0  й=s(��$J"���������Wv�7=`Z��3��c!�~�)��P�-E��fR�е�������c�*���gT9�eQ��-i�C:�}�;�O�u�(ڀI�Ԃ���h��s��q�{F<4���x)D��Ŧm��,t_Ex2D�G�� �_2�=�     