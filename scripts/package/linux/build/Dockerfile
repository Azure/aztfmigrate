FROM mcr.microsoft.com/azurelinux/base/ruby:3.3
ARG version
RUN tdnf install -y gcc gcc-c++ make ruby-devel rpm-build squashfs-tools && \
    tdnf clean all && \
    gem install fpm
CMD ["fpm", "--help"]
