<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:xs="http://www.w3.org/2001/XMLSchema"
    xmlns:s="http://geofon.gfz-potsdam.de/ns/seiscomp3-schema/0.7">

    <!--  
Copyright 2010, Institute of Geological & Nuclear Sciences Ltd or
third-party contributors as indicated by the @author tags.
 
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

    <!-- Converts SeisComPML to a simple CSV format. -->

    <!--Set to the string 'true' to output event info.-->
    <xsl:param name="event"/>
    <!--Set to the string 'true' to output picks info.-->
    <xsl:param name="picks"/>

    <xsl:output method="text"/>

    <xsl:template match="/s:seiscomp">
        <xsl:apply-templates select="s:EventParameters/s:event"/>
    </xsl:template>

    <xsl:template match="s:event">
        <xsl:variable name="preferredOriginID" select="s:preferredOriginID"/>
        <xsl:variable name="preferredMagnitudeID" select="s:preferredMagnitudeID"/>
        <xsl:variable name="publicID" select="@publicID"/>
        <xsl:variable name="type" select="s:type"/>

        <xsl:apply-templates select="../s:origin">
            <xsl:with-param name="preferredOriginID" select="$preferredOriginID"/>
            <xsl:with-param name="preferredMagnitudeID" select="$preferredMagnitudeID"/>
            <xsl:with-param name="publicID" select="$publicID"/>
            <xsl:with-param name="type" select="$type"/>
        </xsl:apply-templates>

    </xsl:template>

    <xsl:template match="s:origin">
        <xsl:param name="preferredOriginID"/>
        <xsl:param name="preferredMagnitudeID"/>
        <xsl:param name="publicID"/>
        <xsl:param name="type"/>

        <xsl:if test="@publicID=$preferredOriginID">

            <xsl:if test="$event='true'">
                <xsl:text>publicid,origintime,latitude,longitude,depth,magnitude,magnitudetype</xsl:text>
                <xsl:value-of select="$newline"/>
                <xsl:value-of select="$publicID"/>
                <xsl:value-of select="$comma"/>
                <xsl:value-of select="s:time/s:value"/>
                <xsl:value-of select="$comma"/>
                <xsl:value-of select="s:latitude/s:value"/>
                <xsl:value-of select="$comma"/>
                <xsl:value-of select="s:longitude/s:value"/>
                <xsl:value-of select="$comma"/>
                <xsl:value-of select="s:depth/s:value"/>

                <xsl:apply-templates select="s:magnitude">
                    <xsl:with-param name="preferredMagnitudeID" select="$preferredMagnitudeID"/>
                </xsl:apply-templates>

                <xsl:value-of select="$newline"/>
            </xsl:if>

            <xsl:if test="$picks='true'">


                <xsl:text>phase,datetime,weight,polarity,network,station,channel,location,mode,status</xsl:text>
                <xsl:value-of select="$newline"/>

                <xsl:for-each select="s:arrival">
                    <xsl:apply-templates select="../../s:pick">
                        <xsl:with-param name="pickID" select="s:pickID"/>
                        <xsl:with-param name="phase" select="s:phase"/>
                        <xsl:with-param name="timeWeight" select="s:weight"/>
                    </xsl:apply-templates>
                </xsl:for-each>
            </xsl:if>
        </xsl:if>
    </xsl:template>


    <xsl:template match="s:pick">
        <xsl:param name="pickID"/>
        <xsl:param name="phase"/>
        <xsl:param name="timeWeight"/>

        <xsl:if test="@publicID=$pickID">
            <xsl:value-of select="$phase"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:time/s:value"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="$timeWeight"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:polarity"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:waveformID/@networkCode"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:waveformID/@stationCode"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:waveformID/@channelCode"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:waveformID/@locationCode"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:evaluationMode"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:evaluationStatus"/>
            <xsl:value-of select="$newline"/>
        </xsl:if>

    </xsl:template>

    <xsl:template match="s:magnitude">
        <xsl:param name="preferredMagnitudeID"/>

        <xsl:if test="@publicID=$preferredMagnitudeID">
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:magnitude/s:value"/>
            <xsl:value-of select="$comma"/>
            <xsl:value-of select="s:type"/>
        </xsl:if>
    </xsl:template>

    <xsl:variable name="newline">
        <xsl:text>&#xa;</xsl:text>
    </xsl:variable>

    <xsl:variable name="comma">
        <xsl:text>,</xsl:text>
    </xsl:variable>

</xsl:stylesheet>
