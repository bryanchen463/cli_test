FROM nginx
COPY conf.d/nginx.conf /etc/nginx/nginx.conf
COPY cert/* /etc/nginx/cert/
EXPOSE 80
EXPOSE 443
CMD ["nginx","-g","daemon off;"]
